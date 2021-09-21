package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	OfferePlansAllPlan      = "ALL_PLAN"
	OfferePlansMedicalPlan  = "MEDICAL_PLAN"
	OfferePlansPensionPlan  = "PENSION_PLAN"
	OfferePlansFamilyPlan   = "FAMILY_PLAN"
	OfferePlansChildrenPlan = "CHILDREN_PLAN"
	OfferePlansTermPlan     = "TERM_PLAN"
	OfferePlansMotorPlan    = "MOTOR_PLAN"
	OfferePlansTravelPlan   = "TRAVEL_PLAN"
)

type Insurance struct {
	Id               uuid.UUID `json:"id" gorm:"type:uuid;PRIMARY_KEY;"`
	Type             string    `json:"type"`
	Premium          float32   `json:"premium"`
	Cover            float32   `json:"cover"`
	YoyDeductionRate float32   `json:"yoy_deduction_rate"`
	gorm.Model
}

type UserExistingInsurance struct {
	Id            uuid.UUID    `json:"id" gorm:"type:uuid;PRIMARY_KEY;"`
	UserConsentId uuid.UUID    `json:"consent_id"`
	Type          string       `json:"type"`
	Premium       float32      `json:"premium"`
	Cover         float32      `json:"cover"`
	AccountId     string       `json:"account_id"`
	UserConsent   UserConsents `gorm:"foreignKey:UserConsentId;constraint:OnDelete:CASCADE;"`
	gorm.Model
}

func (insurance *UserExistingInsurance) BeforeCreate(tx *gorm.DB) (err error) {
	insurance.Id = uuid.New()
	return
}

func (insurance *Insurance) BeforeCreate(tx *gorm.DB) (err error) {
	insurance.Id = uuid.New()
	return
}
