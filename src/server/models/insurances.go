package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
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

var InsuranceNGMockedClauses = []string{
	"We have an SLA of 48 hours of claim resolution. You are entited to 2% addition claim for every additional day.",
	"Your coverage is effective from the time of payment without any delay or waiting period.",
	"Your coverage DOES include all the previously known and declared conditions.",
}

var InsuranceApnaMockedClauses = []string{
	"We have an SLA of 72 hours of claim resolution. You are entited to 1% addition claim for every additional day.",
	"Your coverage is effective from the time of payment without any delay or waiting period.",
	"Your coverage DOES NOT include all the previously known and declared conditions.",
}

type Insurance struct {
	Id               uuid.UUID      `json:"id" gorm:"type:uuid;PRIMARY_KEY;"`
	Title            string         `json:"title"`
	Description      string         `json:"description"`
	Type             string         `json:"type"`
	Premium          float64        `json:"premium"`
	Cover            float64        `json:"cover"`
	Clauses          pq.StringArray `json:"clauses" gorm:"type:text[]"`
	YoyDeductionRate float32        `json:"yoy_deduction_rate"`
	gorm.Model
}

type UserInsurance struct {
	Id                uuid.UUID      `json:"id" gorm:"type:uuid;PRIMARY_KEY;"`
	Type              string         `json:"type"`
	CustomerId        string         `json:"customer_id"`
	Premium           float64        `json:"premium"`
	Cover             float64        `json:"cover"`
	IsActive          bool           `json:"is_active"`
	IsClaimed         bool           `json:"is_claimed"`
	AccountId         string         `json:"account_id"`
	Clauses           pq.StringArray `json:"clauses" gorm:"type:text[]"`
	IsInsuranceNgAcct bool           `json:"is_insurance_ng_acct"`
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
