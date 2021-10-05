package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserScores model is table that stores the scores calculated as per their
// habits and lifestyle based on FIP data.
type UserScores struct {
	ID                uuid.UUID    `json:"id" gorm:"type:uuid;PRIMARY_KEY;"`
	UserConsentID     uuid.UUID    `json:"consent_id"`
	Name              string       `json:"name"`
	DateOfBirth       time.Time    `json:"date_of_birth"`
	Pancard           string       `json:"pan_card"`
	CkycCompliance    bool         `json:"ckyc_compliance"`
	AgeScore          float32      `json:"age_score"`
	WealthScore       float32      `json:"wealth_score"`
	DebtScore         float32      `json:"debt_score"`
	MedicalScore      float32      `json:"medical_score"`
	InvestmentScore   float32      `json:"investment_score"`
	PensionScore      float32      `json:"pension_score"`
	FamilyScore       float32      `json:"family_score"`
	ChildrenScore     float32      `json:"children_score"`
	MotorScore        float32      `json:"motor_score"`
	TermScore         float32      `json:"term_score"`
	TravelScore       float32      `json:"travel_score"`
	AllScore          float32      `json:"all_score"`
	SharedDataSources int16        `json:"shared_data_sources"`
	UserConsent       UserConsents `gorm:"foreignKey:UserConsentID;constraint:OnDelete:CASCADE;"`
}

// BeforeCreate creates new UUID before saving to database
func (userScore *UserScores) BeforeCreate(tx *gorm.DB) (err error) {
	userScore.ID = uuid.New()
	return
}
