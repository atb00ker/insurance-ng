package account_aggregator

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"insurance-ng/src/server/config"
	"insurance-ng/src/server/models"
	"strings"
	"sync"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/google/uuid"
)

func createAndSaveSessionDetails(userId string) (err error) {
	userConsent, err := getUserConsentWithUserId(userId)
	if err != nil {
		return
	}

	rahasyaKeys := getRahasyaKeys()
	sessionData, err := getDataSession(userId, rahasyaKeys, userConsent)
	if err != nil {
		return
	}

	updatedUserConsent := models.UserConsents{
		SessionId:         sessionData.SessionId.String(),
		RahasyaNonce:      rahasyaKeys.KeyMaterial.Nonce,
		RahasyaPrivateKey: rahasyaKeys.PrivateKey,
	}

	config.Database.Where("user_id = ?", userId).Updates(&updatedUserConsent)
	return
}

func saveFipData(sessionId string) (err error) {
	userConsent, err := getUserConsentWithSessionId(sessionId)
	if err != nil {
		return
	}

	// If the data is already fetched, we don't need to do anything.
	if userConsent.DataFetched {
		return
	}

	encryptedData, data_err := getEncryptedFIData(userConsent.SessionId)
	if err != nil {
		err = data_err
		return
	}

	allFipData, err := getDataFromAllFIP(userConsent.RahasyaNonce, userConsent.RahasyaPrivateKey, encryptedData)
	if err != nil {
		return
	}

	if err = processAndSaveFipDataCollection(allFipData, userConsent); err != nil {
		return
	}
	return
}

func getDataSession(userId string, rahasyaKeys rahasyaKeyResponse,
	consentData models.UserConsents) (sessionData setuFiSessionResponse, err error) {
	// TODO: Unfortunatly, there are some bugs in the Setu API as of today, which
	// I have reported, but until those bugs are fixed, we have to comment this and
	// use the hack below.
	// currentTime := time.Now().Format(time.RFC3339)
	uuid := uuid.New()
	currentTime := time.Now()
	currentTimeHack := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.153Z",
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		currentTime.Hour(), currentTime.Minute(), currentTime.Second())

	fiSessionBody := createFiDataRequestBody(uuid, currentTimeHack, consentData, rahasyaKeys)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, fiSessionBody)
	setuRequestBody, err := json.Marshal(fiSessionBody)
	if err != nil {
		return
	}

	respBytes, err := sendRequestToSetu(SetuApiFiRequest, "POST", setuRequestBody, jwtToken)
	if err != nil {
		return
	}

	err = json.Unmarshal(respBytes, &sessionData)
	return
}

func createFiDataRequestBody(uuid uuid.UUID, currentTime string, consentData models.UserConsents,
	rahasyaKeys rahasyaKeyResponse) setuFiSessionRequest {

	signedConsent := strings.Split(consentData.SignedConsent, ".")[2]
	requestBody := setuFiSessionRequest{
		Ver:       "1.0",
		Timestamp: currentTime,
		Txnid:     uuid,
		FIDataRange: fIDataRange{
			From: "1947-08-15T00:00:00.153Z",
			To:   currentTime,
		},
		Consent: fiConsent{
			Id:               consentData.ConsentId,
			DigitalSignature: signedConsent,
		},
		KeyMaterial: rahasyaKeys.KeyMaterial,
	}

	return requestBody
}

func getEncryptedFIData(sessionId string) (fiEncryptedData setuFiDataResponse, err error) {
	urlPath := fmt.Sprintf(SetuApiFiDataFetch, sessionId)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{})
	respBytes, err := sendRequestToSetu(urlPath, "GET", []byte{}, jwtToken)
	if err != nil {
		return
	}

	err = json.Unmarshal(respBytes, &fiEncryptedData)
	return
}

func getDataFromAllFIP(rahasyaNonce string, rahasyaPrivateKey string,
	encryptedData setuFiDataResponse) (response []fipDataCollection, err error) {
	var wgEncrpyptedData sync.WaitGroup

	for _, encryptedFI := range encryptedData.FI {
		wgEncrpyptedData.Add(1)
		go func(encryptedFI fiEncryptionData) {
			defer wgEncrpyptedData.Done()
			fipDataList, err := getFIPData(rahasyaNonce, rahasyaPrivateKey, encryptedFI)
			if err != nil {
				return
			}

			response = append(response, fipDataCollection{
				FipData: fipDataList,
				FipId:   encryptedFI.FipId,
			})
		}(encryptedFI)
	}

	wgEncrpyptedData.Wait()
	return
}

func getFIPData(rahasyaNonce string, rahasyaPrivateKey string,
	encryptedDataList fiEncryptionData) (fipDataList []fipData, err error) {
	var wgEncryptedRecord sync.WaitGroup

	for _, encryptedRecord := range encryptedDataList.Data {
		wgEncryptedRecord.Add(1)
		go func(encryptedRecord fiData) {
			defer wgEncryptedRecord.Done()
			responseData, err := getDecryptedData(rahasyaNonce, rahasyaPrivateKey, encryptedDataList, encryptedRecord)
			if err != nil {
				return
			}

			data, err := base64.StdEncoding.DecodeString(responseData.Base64Data)
			if err != nil {
				return
			}

			var fiData fipData
			if err = xml.Unmarshal(data, &fiData.Account); err != nil {
				return
			}

			fipDataList = append(fipDataList, fiData)
		}(encryptedRecord)
	}

	wgEncryptedRecord.Wait()
	return
}

func getDecryptedData(rahasyaNonce string, rahasyaPrivateKey string,
	encryptedData fiEncryptionData, encryptedRecord fiData) (responseData rahasyaDataResponse, err error) {

	requestBody := rahasyaDecryptRequest{
		Base64Data:        encryptedRecord.EncryptedFI,
		Base64RemoteNonce: encryptedData.KeyMaterial.Nonce,
		Base64YourNonce:   rahasyaNonce,
		OurPrivateKey:     rahasyaPrivateKey,
		RemoteKeyMaterial: encryptedData.KeyMaterial,
	}

	rahasyaRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		return
	}

	respBytes, err := sendRequestToRahasya(RahasyaApiDecrypt, "POST", rahasyaRequestBody)
	if err != nil {
		return
	}

	err = json.Unmarshal(respBytes, &responseData)
	return
}

func getUserData(userId string) (responseData getUserDataResponse, err error) {
	userConsentCh := make(chan userConsentChResp, 1)
	go getUserConsent(userId, userConsentCh)
	insuranceAvailableCh := make(chan insuranceChResp, 1)
	go getAllInsuranceOffers(insuranceAvailableCh)

	insuranceAvailable := <-insuranceAvailableCh
	if insuranceAvailable.err != nil {
		err = insuranceAvailable.err
		return
	}

	userConsent := <-userConsentCh
	if userConsent.err != nil {
		err = userConsent.err
		return
	}

	if !userConsent.result.DataFetched {
		return getUserDataResponse{Status: false}, nil
	}

	userPlanScoresCh := make(chan userPlanScoreChResp, 1)
	go getUserPlanScore(userConsent.result.Id, userPlanScoresCh)
	userExistingInsuranceCh := make(chan userExistingInsurancesChResp, 1)
	go getUserExistingInsurances(userConsent.result.Id, userExistingInsuranceCh)

	userScore := <-userPlanScoresCh
	if userScore.err != nil {
		err = insuranceAvailable.err
		return
	}

	userExistingInsurance := <-userExistingInsuranceCh
	if userExistingInsurance.err != nil {
		err = userConsent.err
		return
	}

	var insuranceOffered []insuranceOffers

	for _, insurance := range insuranceAvailable.result {
		existingInsurance := models.UserExistingInsurance{}
		for _, userInsurances := range userExistingInsurance.result {
			if userInsurances.Type == insurance.Type {
				existingInsurance = *userInsurances
				break
			}
		}

		insuranceOffered = append(insuranceOffered, insuranceOffers{
			Title:          insurance.Title,
			Description:    insurance.Description,
			AccountId:      existingInsurance.AccountId,
			Score:          getScoreForType(userScore.result, insurance.Type),
			CurrentPremium: existingInsurance.Premium,
			CurrentCover:   existingInsurance.Cover,
			OfferedPremium: insurance.Premium,
			OfferedCover:   insurance.Cover,
			Type:           insurance.Type,
		})
	}

	responseData = getUserDataResponse{
		Status: true,
		Data: userData{
			Name:            userScore.result.Name,
			DateOfBirth:     userScore.result.DateOfBirth,
			Pancard:         userScore.result.PanCard,
			CkycCompliance:  userScore.result.CkycCompliance,
			AgeScore:        userScore.result.AgeScore,
			WealthScore:     userScore.result.WealthScore,
			DebtScore:       userScore.result.DebtScore,
			InvestmentScore: userScore.result.InvestmentScore,
			InsuranceOffers: insuranceOffered,
		},
		Error: nil,
	}

	return
}

func getScoreForType(userScore *models.UserPlanScores, insuranceType string) float32 {
	switch insuranceType {
	case models.OfferePlansAllPlan:
		return userScore.AllScore
	case models.OfferePlansMedicalPlan:
		return userScore.MedicalScore
	case models.OfferePlansPensionPlan:
		return userScore.PensionScore
	case models.OfferePlansFamilyPlan:
		return userScore.FamilyScore
	case models.OfferePlansChildrenPlan:
		return userScore.ChildrenScore
	case models.OfferePlansTermPlan:
		return userScore.TermScore
	case models.OfferePlansMotorPlan:
		return userScore.MotorScore
	case models.OfferePlansTravelPlan:
		return userScore.TravelScore
	}
	return 1
}

func getUserExistingInsurances(consentId uuid.UUID, userExistingInsuranceCh chan userExistingInsurancesChResp) {
	var userExistingInsurance []*models.UserExistingInsurance
	result := config.Database.Where("user_consent_id = ?", consentId).Find(&userExistingInsurance)

	userExistingInsuranceCh <- userExistingInsurancesChResp{
		result: userExistingInsurance,
		err:    result.Error,
	}
}

func getUserPlanScore(consentId uuid.UUID, userPlanScoreCh chan userPlanScoreChResp) {
	var userPlanScore *models.UserPlanScores
	result := config.Database.Where("user_consent_id = ?", consentId).Take(&userPlanScore)

	userPlanScoreCh <- userPlanScoreChResp{
		result: userPlanScore,
		err:    result.Error,
	}
}

func getUserConsent(userId string, userConsentChannel chan userConsentChResp) {
	var userConsent *models.UserConsents
	result := config.Database.Where("user_id = ?", userId).Take(&userConsent)

	userConsentChannel <- userConsentChResp{
		result: userConsent,
		err:    result.Error,
	}
}

func getAllInsuranceOffers(insuranceAvailableChannel chan insuranceChResp) {
	var insurance []*models.Insurance
	result := config.Database.Find(&insurance)

	insuranceAvailableChannel <- insuranceChResp{
		result: insurance,
		err:    result.Error,
	}
}
