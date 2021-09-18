package account_aggregator

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"insurance-ng/src/server/config"
	"insurance-ng/src/server/models"
	"strings"
	"sync"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/google/uuid"
)

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

func getEncryptedFIData(sessionData setuFiSessionResponse) (fiEncryptedData setuFiDataResponse, err error) {
	urlPath := fmt.Sprintf(SetuApiFiDataFetch, sessionData.SessionId.String())
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{})
	respBytes, err := sendRequestToSetu(urlPath, "GET", []byte{}, jwtToken)
	if err != nil {
		return
	}

	err = json.Unmarshal(respBytes, &fiEncryptedData)
	return
}

func getUnencryptedFIDataList(rahasyaKeys rahasyaKeyResponse,
	encryptedData setuFiDataResponse) (response []rahasyaDataResponseCollection, err error) {
	var wgEncrpyptedData sync.WaitGroup

	for _, encryptedFI := range encryptedData.FI {
		wgEncrpyptedData.Add(1)
		go func(encryptedFI fiEncryptionData) {
			defer wgEncrpyptedData.Done()
			fiData, err := prepareFIForDecryption(rahasyaKeys, encryptedFI)
			if err != nil {
				return
			}

			fiDataCollection := rahasyaDataResponseCollection{
				RahasyaData: fiData,
				FipId:       encryptedFI.FipId,
			}

			response = append(response, fiDataCollection)
		}(encryptedFI)
	}

	wgEncrpyptedData.Wait()
	return
}

func prepareFIForDecryption(rahasyaKeys rahasyaKeyResponse,
	encryptedDataList fiEncryptionData) (response []rahasyaDataResponse, err error) {
	var wgEncrpyptedData sync.WaitGroup

	for _, encryptedRecord := range encryptedDataList.Data {
		wgEncrpyptedData.Add(1)
		go func(encryptedRecord fiData) {
			defer wgEncrpyptedData.Done()
			responseData, err := getDecryptedData(rahasyaKeys, encryptedDataList, encryptedRecord)
			if err != nil {
				response = append(response, rahasyaDataResponse{
					Data:      "",
					ErrorInfo: err.Error(),
				})
			}

			data, err := base64.StdEncoding.DecodeString(responseData.Base64Data)
			if err != nil {
				response = append(response, rahasyaDataResponse{
					Data:      "",
					ErrorInfo: err.Error(),
				})
			}

			response = append(response, rahasyaDataResponse{
				Data:      string(data),
				ErrorInfo: responseData.ErrorInfo,
			})
		}(encryptedRecord)
	}

	wgEncrpyptedData.Wait()
	return
}

func getDecryptedData(rahasyaKeys rahasyaKeyResponse,
	encryptedData fiEncryptionData, encryptedRecord fiData) (responseData rahasyaDataResponse, err error) {

	requestBody := rahasyaDecryptRequest{
		Base64Data:        encryptedRecord.EncryptedFI,
		Base64RemoteNonce: encryptedData.KeyMaterial.Nonce,
		Base64YourNonce:   rahasyaKeys.KeyMaterial.Nonce,
		OurPrivateKey:     rahasyaKeys.PrivateKey,
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

func savebase64Data(base64Data []rahasyaDataResponseCollection, userId string) error {
	userData, err := json.Marshal(base64Data)
	if err != nil {
		return err
	}
	updatedRow := config.Database.Model(&models.UserConsents{}).Where(
		"user_id = ?", userId).Update("user_data", userData)

	return updatedRow.Error
}
