package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPlanScores struct {
	Id             uuid.UUID    `json:"id" gorm:"type:uuid;PRIMARY_KEY;"`
	UserConsentId  uuid.UUID    `json:"consent_id"`
	Name           string       `json:"name"`
	DateOfBirth    time.Time    `json:"date_of_birth"`
	PanCard        string       `json:"pan_card"`
	CkycCompliance bool         `json:"ckyc_compliance"`
	AgeScore       float32      `json:"age_score"`
	MedicalScore   float32      `json:"medical_score"`
	PensionScore   float32      `json:"pension_score"`
	FamilyScore    float32      `json:"family_score"`
	ChildrenScore  float32      `json:"children_score"`
	MotorScore     float32      `json:"motor_score"`
	TermScore      float32      `json:"term_score"`
	TravelScore    float32      `json:"travel_score"`
	AllScore       float32      `json:"all_score"`
	UserConsent    UserConsents `gorm:"foreignKey:UserConsentId;constraint:OnDelete:CASCADE;"`
}

func (insurance *UserPlanScores) BeforeCreate(tx *gorm.DB) (err error) {
	insurance.Id = uuid.New()
	return
}
