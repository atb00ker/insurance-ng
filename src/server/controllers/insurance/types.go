package insurance

import (
	"insurance-ng/src/server/models"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/lib/pq"
)

const (
	PreApprovedBar = 0.80
)

type purchaseRequest struct {
	Uuid uuid.UUID `json:"uuid"`
}

// Get User Data Types //

type getUserDataResponse struct {
	Status bool     `json:"status"`
	Data   userData `json:"data"`
	Error  error    `json:"error"`
}

type userData struct {
	Name              string            `json:"name"`
	DateOfBirth       time.Time         `json:"date_of_birth"`
	Pancard           string            `json:"pancard"`
	Phone             string            `json:"phone"`
	CkycCompliance    bool              `json:"ckyc_compliance"`
	AgeScore          float32           `json:"age_score"`
	WealthScore       float32           `json:"wealth_score"`
	DebtScore         float32           `json:"debt_score"`
	InvestmentScore   float32           `json:"investment_score"`
	InsuranceOffers   []insuranceOffers `json:"insurance"`
	SharedDataSources int16             `json:"shared_data_sources"`
}

type insuranceOffers struct {
	Id                uuid.UUID      `json:"uuid"`
	AccountId         string         `json:"account_id"`
	Type              string         `json:"type"`
	Title             string         `json:"title"`
	Description       string         `json:"description"`
	Score             float32        `json:"score"`
	CurrentPremium    float64        `json:"current_premium"`
	CurrentCover      float64        `json:"current_cover"`
	OfferedPremium    float64        `json:"offer_premium"`
	OfferedCover      float64        `json:"offer_cover"`
	YoyDeductionRate  float32        `json:"yoy_deduction_rate"`
	Clauses           pq.StringArray `json:"clauses"`
	CurrentClauses    pq.StringArray `json:"current_clauses"`
	IsInsuranceNgAcct bool           `json:"is_insurance_ng_acct"`
	IsActive          bool           `json:"is_active"`
}

type userScoreChResp struct {
	result *models.UserScores
	err    error
}
type userInsurancesChResp struct {
	result []*models.UserInsurance
	err    error
}

type userConsentChResp struct {
	result *models.UserConsents
	err    error
}

type insuranceChResp struct {
	result []*models.Insurance
	err    error
}

// Dashboard Websocket //

type Websocket struct {
	clients    map[*string]*Client
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	Id         *string
	websocket  *Websocket
	connection *websocket.Conn
}
