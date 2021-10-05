package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Plans offered by Insurance NG
const (
	OfferPlansAllPlan      = "ALL_PLAN"
	OfferPlansMedicalPlan  = "MEDICAL_PLAN"
	OfferPlansPensionPlan  = "PENSION_PLAN"
	OfferPlansFamilyPlan   = "FAMILY_PLAN"
	OfferPlansChildrenPlan = "CHILDREN_PLAN"
	OfferPlansTermPlan     = "TERM_PLAN"
	OfferPlansMotorPlan    = "MOTOR_PLAN"
	OfferPlansTravelPlan   = "TRAVEL_PLAN"
)

// InsuranceNGMockedClauses is mocked insurance clauses of other insurance providers
var InsuranceNGMockedClauses = []string{
	"We have an SLA of 48 hours of claim resolution. You are entited to 2% addition claim for every additional day.",
	"Your coverage is effective from the time of payment without any delay or waiting period.",
	"Your coverage DOES include all the previously known and declared conditions.",
}

// InsuranceApnaMockedClauses are clauses offered by Insurance NG
var InsuranceApnaMockedClauses = []string{
	"We have an SLA of 72 hours of claim resolution. You are entited to 1% addition claim for every additional day.",
	"Your coverage is effective from the time of payment without any delay or waiting period.",
	"Your coverage DOES NOT include all the previously known and declared conditions.",
}

// Insurance model is table of all Insurances offered by Insurance NG.
type Insurance struct {
	ID               uuid.UUID      `json:"id" gorm:"type:uuid;PRIMARY_KEY;"`
	Title            string         `json:"title"`
	Description      string         `json:"description"`
	Type             string         `json:"type"`
	Premium          float64        `json:"premium"`
	Cover            float64        `json:"cover"`
	Clauses          pq.StringArray `json:"clauses" gorm:"type:text[]"`
	YoyDeductionRate float32        `json:"yoy_deduction_rate"`
	gorm.Model
}

// UserInsurance model is table that stores information about insurances that are purchased by the user
type UserInsurance struct {
	ID                uuid.UUID      `json:"id" gorm:"type:uuid;PRIMARY_KEY;"`
	Type              string         `json:"type"`
	CustomerID        string         `json:"customer_id"`
	Premium           float64        `json:"premium"`
	Cover             float64        `json:"cover"`
	IsActive          bool           `json:"is_active"`
	IsClaimed         bool           `json:"is_claimed"`
	AccountID         string         `json:"account_id"`
	Clauses           pq.StringArray `json:"clauses" gorm:"type:text[]"`
	IsInsuranceNgAcct bool           `json:"is_insurance_ng_acct"`
	gorm.Model
}

// BeforeCreate creates new UUID before saving to database
func (userInsurance *UserInsurance) BeforeCreate(tx *gorm.DB) (err error) {
	userInsurance.ID = uuid.New()
	return
}

// BeforeCreate creates new UUID before saving to database
func (insurance *Insurance) BeforeCreate(tx *gorm.DB) (err error) {
	insurance.ID = uuid.New()
	return
}
