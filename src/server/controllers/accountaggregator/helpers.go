package accountaggregator

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"insurance-ng/src/server/config"
	"insurance-ng/src/server/models"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/google/uuid"
)

func getDetachedJwt(jwtToken *jwt.Token) (string, error) {
	privateKey, err := getPrivatePemFileContent()
	if err != nil {
		return "", err
	}

	signedJwtToken, err := jwtToken.SignedString(privateKey)

	if err != nil {
		return "", err
	}

	jwtTokenParts := strings.Split(signedJwtToken, ".")
	jwtTokenParts[1] = ""
	return strings.Join(jwtTokenParts, "."), nil
}

func sendRequestToSetu(urlPath string, reqType string, payload []byte,
	jwtToken *jwt.Token) (response []byte, err error) {

	setuAPI, err := url.Parse(os.Getenv("APP_SETU_AA_ENDPOINT"))
	if err != nil {
		return
	}
	setuAPI.Path = path.Join(setuAPI.Path, urlPath)

	setuRequest, err := http.NewRequest(reqType, setuAPI.String(), bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	requestJWS, err := getDetachedJwt(jwtToken)
	if err != nil {
		return
	}

	clientAPI := os.Getenv("APP_SETU_CLIENT_KEY")
	setuRequest.Header = http.Header{
		"Content-Type":    []string{"application/json"},
		"client_api_key":  []string{clientAPI},
		"x-jws-signature": []string{requestJWS},
	}

	client := &http.Client{}
	setuResponse, err := client.Do(setuRequest)
	if err != nil {
		return
	}
	defer setuResponse.Body.Close()

	response, err = ioutil.ReadAll(setuResponse.Body)
	return
}

func sendResponseToSetuNotification() (clientAPI string, requestJWS string,
	setuResponseBody []byte, err error) {

	// Time Hack:
	// Currently, Setu API is not following the RFC3339 Correctly,
	// Hence for the time being, we manually converting dates.
	startTime := time.Now()
	startTime3339 := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.153Z",
		startTime.Year(), startTime.Month(), startTime.Day(),
		startTime.Hour(), startTime.Minute(), startTime.Second())

	respMessage := setuConsentNotificationResponse{
		Ver:       "1.0",
		Timestamp: startTime3339,
		Txnid:     uuid.New(),
		Response:  "OK",
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, respMessage)
	setuResponseBody, err = json.Marshal(respMessage)
	if err != nil {
		return
	}

	requestJWS, err = getDetachedJwt(jwtToken)
	if err != nil {
		return
	}

	clientAPI = os.Getenv("APP_SETU_CLIENT_KEY")
	return
}

func handleNotificationError(response http.ResponseWriter, err error) {
	// Time Hack:
	// Currently, Setu API is not following the RFC3339 Correctly,
	// Hence for the time being, we manually converting dates.
	startTime := time.Now()
	startTime3339 := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.153Z",
		startTime.Year(), startTime.Month(), startTime.Day(),
		startTime.Hour(), startTime.Minute(), startTime.Second())

	response.WriteHeader(http.StatusBadRequest)
	respMessage, _ := json.Marshal(setuConsentNotificationResponse{
		ErrorCode: "Errored",
		Ver:       "1.0",
		Timestamp: startTime3339,
		Txnid:     uuid.New(),
		Response:  err.Error(),
	})
	response.Write(respMessage)
}

func getUserConsentWithUserID(userID string) (userConsent models.UserConsents, err error) {
	result := config.Database.Model(&models.UserConsents{}).Where("user_id = ?", userID).Take(&userConsent)
	err = result.Error
	return
}

func getUserConsentWithSessionID(sessionID string) (userConsent models.UserConsents, err error) {
	result := config.Database.Model(&models.UserConsents{}).Where("session_id = ?", sessionID).Take(&userConsent)
	err = result.Error
	return
}

func getPrivatePemFileContent() (x509Key interface{}, err error) {
	filePath := os.Getenv("APP_SETU_JWS_PRIVATEKEY_PATH")
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println(fmt.Sprintf("Filepath (%s) does not point to a readable file, please check the value of `APP_SETU_JWS_PRIVATEKEY_PATH`.", filePath))
	}

	block, _ := pem.Decode(content)
	// Check if it's a private key
	if block == nil || block.Type != "PRIVATE KEY" {
		err = errors.New("failed to decode PEM block containing private key")
		return
	}

	x509Key, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	return
}

// Rahasya //

func getRahasyaKeys() (rahasyaKeys rahasyaKeyResponse, err error) {
	respBytes, err := sendRequestToRahasya(RahasyaAPIGetKeys, "GET", []byte{})
	if err != nil {
		return
	}

	json.Unmarshal(respBytes, &rahasyaKeys)
	return
}

func deleteUserConsent(userConsent models.UserConsents) (err error) {
	// We delete here instead of updating because we want to delete
	// cascade all FIP the stored with this user consent.
	if result := config.Database.Where("customer_id = ?", userConsent.CustomerID).Where(
		"is_insurance_ng_acct <> ?", true).Delete(&models.UserInsurance{}); result.Error != nil {
		return result.Error
	}

	if result := config.Database.Where("user_id = ?",
		userConsent.UserID).Delete(&userConsent); result.Error != nil {
		return result.Error
	}
	return
}

func sendRequestToRahasya(urlPath string, reqType string, payload []byte) (response []byte, err error) {
	rahasyaAPI, err := url.Parse(os.Getenv("APP_RAHASYA_AA_ENDPOINT"))
	if err != nil {
		return
	}
	rahasyaAPI.Path = path.Join(rahasyaAPI.Path, urlPath)

	rahasyaRequest, err := http.NewRequest(reqType, rahasyaAPI.String(), bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	rahasyaRequest.Header = http.Header{
		"Content-Type": []string{"application/json"},
	}

	client := &http.Client{}
	rahasyaResponse, err := client.Do(rahasyaRequest)
	if err != nil {
		return
	}
	defer rahasyaResponse.Body.Close()

	response, err = ioutil.ReadAll(rahasyaResponse.Body)
	return
}
