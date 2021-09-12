package account_aggregator

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"insurance-ng/src/server/config"
	"insurance-ng/src/server/models"
	"strings"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/google/uuid"
)

func getDataSession(userId string, rahasyaKeys rahasyaKeyResponse,
	consentData models.UserConsents) (sessionData fiSessionResponse, err error) {
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
	rahasyaKeys rahasyaKeyResponse) fiSessionRequest {

	signedConsent := strings.Split(consentData.SignedConsent, ".")[2]
	requestBody := fiSessionRequest{
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

func getEncryptedFIData(sessionData fiSessionResponse) (fiEncryptedData fiDataResponse, err error) {
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
	encryptedData fiDataResponse) (response []rahasyaDataResponseCollection, err error) {

	for _, encryptedFI := range encryptedData.FI {
		// TODO: Make Async
		var fiData []rahasyaDataResponse
		fiData, err = prepareFIForDecryption(rahasyaKeys, encryptedFI)
		if err != nil {
			return
		}

		fiDataCollection := rahasyaDataResponseCollection{
			RahasyaData: fiData,
			FipId:       encryptedFI.FipId,
		}
		response = append(response, fiDataCollection)
	}

	return
}

func prepareFIForDecryption(rahasyaKeys rahasyaKeyResponse,
	encryptedData fiEncryptionData) (response []rahasyaDataResponse, err error) {

	for _, encryptedRecord := range encryptedData.Data {
		// TODO: Make Async
		var responseData rahasyaDataResponse
		responseData, err = sendDecryptRequestToRahasya(rahasyaKeys, encryptedData, encryptedRecord)
		if err != nil {
			response = append(response, rahasyaDataResponse{
				Data:      "",
				ErrorInfo: err.Error(),
			})
		}

		var data []byte
		data, err = base64.StdEncoding.DecodeString(responseData.Base64Data)
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
	}

	return
}

func sendDecryptRequestToRahasya(rahasyaKeys rahasyaKeyResponse,
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
