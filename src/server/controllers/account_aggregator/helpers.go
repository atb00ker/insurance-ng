package account_aggregator

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

	setuApi, err := url.Parse(os.Getenv("APP_SETU_AA_ENDPOINT"))
	if err != nil {
		return
	}
	setuApi.Path = path.Join(setuApi.Path, urlPath)

	setuRequest, err := http.NewRequest(reqType, setuApi.String(), bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	requestJws, err := getDetachedJwt(jwtToken)
	if err != nil {
		return
	}

	clientApi := os.Getenv("APP_SETU_CLIENT_KEY")
	setuRequest.Header = http.Header{
		"Content-Type":    []string{"application/json"},
		"client_api_key":  []string{clientApi},
		"x-jws-signature": []string{requestJws},
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

func sendResponseToSetuNotification() (clientApi string, requestJws string,
	setuResponseBody []byte, err error) {

	// Hack
	// Currently, Setu API is not following
	// the complete RFC3339
	startTime := time.Now()
	startTimeHack := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.153Z",
		startTime.Year(), startTime.Month(), startTime.Day(),
		startTime.Hour(), startTime.Minute(), startTime.Second())

	respMessage := setuConsentNotificationResponse{
		Ver:       "1.0",
		Timestamp: startTimeHack,
		Txnid:     uuid.New(),
		Response:  "OK",
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, respMessage)
	setuResponseBody, err = json.Marshal(respMessage)
	if err != nil {
		return
	}

	requestJws, err = getDetachedJwt(jwtToken)
	if err != nil {
		return
	}

	clientApi = os.Getenv("APP_SETU_CLIENT_KEY")
	return
}

func HandleNotificationError(response http.ResponseWriter, err error) {
	// Hack
	// Currently, Setu API is not following
	// the complete RFC3339
	startTime := time.Now()
	startTimeHack := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.153Z",
		startTime.Year(), startTime.Month(), startTime.Day(),
		startTime.Hour(), startTime.Minute(), startTime.Second())

	response.WriteHeader(http.StatusBadRequest)
	respMessage, _ := json.Marshal(setuConsentNotificationResponse{
		ErrorCode: "Errored",
		Ver:       "1.0",
		Timestamp: startTimeHack,
		Txnid:     uuid.New(),
		Response:  err.Error(),
	})
	response.Write(respMessage)
}

func getUserConsentWithUserId(userId string) (userConsent models.UserConsents, err error) {
	result := config.Database.Model(&models.UserConsents{}).Where("user_id = ?", userId).Take(&userConsent)
	err = result.Error
	return
}

func getUserConsentWithSessionId(sessionId string) (userConsent models.UserConsents, err error) {
	result := config.Database.Model(&models.UserConsents{}).Where("session_id = ?", sessionId).Take(&userConsent)
	err = result.Error
	return
}

func getPrivatePemFileContent() (x509Key interface{}, err error) {
	filePath := os.Getenv("APP_JWS_AA_PRIVATEKEY_PATH")
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println(fmt.Sprintf("Filepath (%s) does not point to a readable file, please check the value of `APP_JWS_AA_PRIVATEKEY_PATH`.", filePath))
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

func getRahasyaKeys() (rahasyaKeys rahasyaKeyResponse) {
	respBytes, err := sendRequestToRahasya(RahasyaApiGetKeys, "GET", []byte{})
	if err != nil {
		return
	}

	json.Unmarshal(respBytes, &rahasyaKeys)
	return
}

func deleteUserConsent(userConsent models.UserConsents) {
	// We delete here instead of updating because we want to delete
	// cascade all FIP the stored with this user consent.
	var userScore models.UserScores
	config.Database.Model(&models.UserScores{}).Where("user_consent_id = ?",
		userConsent.Id).Take(&userScore)

	config.Database.Where("pancard = ?", userScore.Pancard).Where(
		"is_insuranceng_account <> ?", 1).Delete(&models.UserInsurance{})

	config.Database.Where("user_id = ?", userConsent.UserId).Delete(&userConsent)
}

func sendRequestToRahasya(urlPath string, reqType string, payload []byte) (response []byte, err error) {
	rahasyaApi, err := url.Parse(os.Getenv("APP_RAHASYA_AA_ENDPOINT"))
	if err != nil {
		return
	}
	rahasyaApi.Path = path.Join(rahasyaApi.Path, urlPath)

	rahasyaRequest, err := http.NewRequest(reqType, rahasyaApi.String(), bytes.NewBuffer(payload))
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
