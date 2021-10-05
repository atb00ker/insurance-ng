package insurance

import (
	"encoding/json"
	"insurance-ng/src/server/models"
	"io"
	"math"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func getScoreForType(userScore *models.UserScores, insuranceType string) float32 {
	switch insuranceType {
	case models.OfferPlansAllPlan:
		return userScore.AllScore
	case models.OfferPlansMedicalPlan:
		return userScore.MedicalScore
	case models.OfferPlansPensionPlan:
		return userScore.PensionScore
	case models.OfferPlansFamilyPlan:
		return userScore.FamilyScore
	case models.OfferPlansChildrenPlan:
		return userScore.ChildrenScore
	case models.OfferPlansTermPlan:
		return userScore.TermScore
	case models.OfferPlansMotorPlan:
		return userScore.MotorScore
	case models.OfferPlansTravelPlan:
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

func getInsuranceUUID(requestBody io.Reader) (insuranceUUID uuid.UUID, err error) {
	var requestJSON insuranceActionRequest
	decoder := json.NewDecoder(requestBody)
	if err = decoder.Decode(&requestJSON); err != nil {
		return
	}
	return requestJSON.UUID, nil
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
