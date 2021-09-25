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
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Type             string    `json:"type"`
	Premium          float64   `json:"premium"`
	Cover            float64   `json:"cover"`
	YoyDeductionRate float32   `json:"yoy_deduction_rate"`
	gorm.Model
}

type UserInsurance struct {
	Id                uuid.UUID    `json:"id" gorm:"type:uuid;PRIMARY_KEY;"`
	UserConsentId     uuid.UUID    `json:"consent_id"`
	Type              string       `json:"type"`
	Premium           float64      `json:"premium"`
	Cover             float64      `json:"cover"`
	IsActive          bool         `json:"is_active"`
	AccountId         string       `json:"account_id"`
	IsInsuranceNgAcct bool         `json:"is_insuranceng_account"`
	UserConsent       UserConsents `gorm:"foreignKey:UserConsentId;constraint:OnDelete:CASCADE;"`
	gorm.Model
}

func (insurance *UserInsurance) BeforeCreate(tx *gorm.DB) (err error) {
	insurance.Id = uuid.New()
	return
}

func (insurance *Insurance) BeforeCreate(tx *gorm.DB) (err error) {
	insurance.Id = uuid.New()
	return
}
