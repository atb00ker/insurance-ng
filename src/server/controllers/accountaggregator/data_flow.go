package accountaggregator

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"insurance-ng/src/server/config"
	"insurance-ng/src/server/models"
	"strings"
	"sync"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/google/uuid"
)

func createAndSaveSessionDetails(userID string) (err error) {
	userConsent, err := getUserConsentWithUserID(userID)
	if err != nil {
		return
	}

	rahasyaKeys, err := getRahasyaKeys()
	if err != nil {
		return
	}

	sessionData, err := getDataSession(userID, rahasyaKeys, userConsent)
	if err != nil {
		return
	}

	updatedUserConsent := models.UserConsents{
		SessionID:         sessionData.SessionID.String(),
		RahasyaNonce:      rahasyaKeys.KeyMaterial.Nonce,
		RahasyaPrivateKey: rahasyaKeys.PrivateKey,
	}

	config.Database.Where("user_id = ?", userID).Updates(&updatedUserConsent)
	return
}

func saveFipData(sessionID string) (err error) {
	userConsent, err := getUserConsentWithSessionID(sessionID)
	if err != nil {
		return
	}

	// If the data is already fetched, we don't need to do anything.
	if userConsent.DataFetched {
		return
	}

	encryptedData, dataErr := getEncryptedFIData(userConsent.SessionID)
	if err != nil {
		err = dataErr
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

func getDataSession(userID string, rahasyaKeys rahasyaKeyResponse,
	consentData models.UserConsents) (sessionData setuFiSessionResponse, err error) {
	// Time Hack:
	// Currently, Setu API is not following the RFC3339 Correctly,
	// Hence for the time being, we manually converting dates.
	uuid := uuid.New()
	currentTime := time.Now()
	currentTime3339 := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.153Z",
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		currentTime.Hour(), currentTime.Minute(), currentTime.Second())

	fiSessionBody := createFiDataRequestBody(uuid, currentTime3339, consentData, rahasyaKeys)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, fiSessionBody)
	setuRequestBody, err := json.Marshal(fiSessionBody)
	if err != nil {
		return
	}

	respBytes, err := sendRequestToSetu(SetuAPIFiRequest, "POST", setuRequestBody, jwtToken)
	if err != nil {
		return
	}

	err = json.Unmarshal(respBytes, &sessionData)
	if sessionData.ErrorCode != "" {
		return sessionData, errors.New(sessionData.ErrorMsg)
	}

	return
}

func createFiDataRequestBody(uuid uuid.UUID, currentTime string, consentData models.UserConsents,
	rahasyaKeys rahasyaKeyResponse) setuFiSessionRequest {
	var signedConsent string
	signedConsentList := strings.Split(consentData.SignedConsent, ".")
	if len(signedConsentList) > 2 {
		signedConsent = signedConsentList[2]
	} else {
		return setuFiSessionRequest{}
	}
	requestBody := setuFiSessionRequest{
		Ver:       "1.0",
		Timestamp: currentTime,
		Txnid:     uuid,
		FIDataRange: fIDataRange{
			From: "1947-08-15T00:00:00.153Z",
			To:   currentTime,
		},
		Consent: fiConsent{
			ID:               consentData.ConsentID,
			DigitalSignature: signedConsent,
		},
		KeyMaterial: rahasyaKeys.KeyMaterial,
	}

	return requestBody
}

func getEncryptedFIData(sessionID string) (fiEncryptedData setuFiDataResponse, err error) {
	urlPath := fmt.Sprintf(SetuAPIFiDataFetch, sessionID)
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
				fmt.Println(err.Error())
				return
			}

			response = append(response, fipDataCollection{
				FipData: fipDataList,
				FipID:   encryptedFI.FipID,
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
			if err = json.Unmarshal(data, &fiData); err != nil {
				fmt.Println(err.Error())
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

	respBytes, err := sendRequestToRahasya(RahasyaAPIDecrypt, "POST", rahasyaRequestBody)
	if err != nil {
		return
	}

	err = json.Unmarshal(respBytes, &responseData)
	return
}
