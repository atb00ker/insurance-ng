package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	OfferePlansAllPlan      = "ALL_PLAN"
	OfferePlansMedicalPlan  = "MEDICAL_PLAN"
	OfferePlansPensionPlan  = "PENSION_PLAN"
	OfferePlansHomePlan     = "HOME_PLAN"
	OfferePlansFamilyPlan   = "FAMILY_PLAN"
	OfferePlansChildrenPlan = "CHILDREN_PLAN"
	OfferePlansLifePlan     = "LIFE_PLAN"
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

func (insurance *Insurance) BeforeCreate(tx *gorm.DB) (err error) {
	insurance.Id = uuid.New()
	return
}
