package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ConsentPending  = "PENDING"
	ConsentReady    = "READY"
	ConsentNotFound = "NOTFOUND"
	ConsentError    = "UNKNOWN"
)

const (
	SignedConsentAccepted = "ACTIVE"
	SignedConsentRejected = "REJECTED"
	SignedConsentRevoked  = "REVOKED"
	SignedConsentNotFound = "NOTFOUND"
	SignedConsentError    = "UNKNOWN"
)

type UserConsents struct {
	Id            uuid.UUID `json:"id" gorm:"type:uuid;PRIMARY_KEY;"`
	UserId        string    `json:"user_id"`
	CustomerId    string    `json:"customer_id"`
	Expire        time.Time `json:"expire"`
	Status        string    `json:"status"`
	UserData      string    `json:"user_data" gorm:"default:-"`
	SignedConsent string    `json:"signed_consent" gorm:"default:-"`
	ConsentHandle uuid.UUID `json:"consent_handle" gorm:"type:uuid;"`
	ConsentId     uuid.UUID `json:"consent_id" gorm:"type:uuid;"`
	User          Users     `gorm:"foreignKey:UserId"`
	gorm.Model
}

func (consent *UserConsents) BeforeCreate(tx *gorm.DB) (err error) {
	consent.Id = uuid.New()
	return
}
