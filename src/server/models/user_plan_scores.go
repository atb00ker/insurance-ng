package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPlanScores struct {
	Id                uuid.UUID    `json:"id" gorm:"type:uuid;PRIMARY_KEY;"`
	UserConsentId     uuid.UUID    `json:"consent_id"`
	Name              string       `json:"name"`
	DateOfBirth       time.Time    `json:"date_of_birth"`
	PanCard           string       `json:"pan_card"`
	CkycCompliance    bool         `json:"ckyc_compliance"`
	AgeScore          float32      `json:"age_score"`
	MedicalScore      float32      `json:"medical_score"`
	PensionScore      float32      `json:"pension_score"`
	FamilyScore       float32      `json:"family_score"`
	ChildrenScore     float32      `json:"children_score"`
	MotorScore        float32      `json:"motor_score"`
	TermScore         float32      `json:"term_score"`
	TravelScore       float32      `json:"travel_score"`
	AllScore          float32      `json:"all_score"`
	MedicalAccountId  string       `json:"medical_account_id"`
	PensionAccountId  string       `json:"pension_account_id"`
	FamilyAccountId   string       `json:"family_account_id"`
	ChildrenAccountId string       `json:"children_account_id"`
	MotorAccountId    string       `json:"motor_account_id"`
	TermAccountId     string       `json:"term_account_id"`
	TravelAccountId   string       `json:"travel_account_id"`
	AllAccountId      string       `json:"all_account_id"`
	UserConsent       UserConsents `gorm:"foreignKey:UserConsentId;constraint:OnDelete:CASCADE;"`
}

func (insurance *UserPlanScores) BeforeCreate(tx *gorm.DB) (err error) {
	insurance.Id = uuid.New()
	return
}
