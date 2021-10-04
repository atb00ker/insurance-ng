package insurance

import (
	"encoding/json"
	"insurance-ng/src/server/models"
	"io"
	"math"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func GetScoreForType(userScore *models.UserScores, insuranceType string) float32 {
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

func getOfferPremium(premium float64, score float64) float64 {
	return math.Ceil(premium - ((premium * score) / 15))
}

func getOfferCover(cover float64, score float64) float64 {
	return math.Floor(cover + ((cover * score) / 15))
}

func getInsuranceUuid(requestBody io.Reader) (insuranceUuid uuid.UUID, err error) {
	var requestJson insuranceActionRequest
	decoder := json.NewDecoder(requestBody)
	if err = decoder.Decode(&requestJson); err != nil {
		return
	}
	return requestJson.Uuid, nil
}

func websocketResponse(client *Client, message []byte) (err error) {
	response, err := client.connection.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}

	response.Write(message)
	if err = response.Close(); err != nil {
		return
	}

	return nil
}
