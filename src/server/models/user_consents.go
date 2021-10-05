package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Specical column states
const (
	EmptyColumn       = "-"
	StatusDataFetched = "FETCHED"
)

// Possible states of Artefact Request
const (
	ArtefactStatusReady   = "READY"
	ArtefactStatusActive  = "ACTIVE"
	ArtefactStatusDenied  = "DENIED"
	ArtefactStatusPending = "PENDING"
	ArtefactStatusTimeout = "TIMEOUT"
	ArtefactStatusError   = "UNKNOWN"
)

// Possible states of Consent Status
const (
	ConsentStatusActive   = "ACTIVE"
	ConsentStatusRejected = "REJECTED"
	ConsentStatusRevoked  = "REVOKED"
	ConsentStatusPaused   = "PAUSED"
	ConsentStatusError    = "UNKNOWN"
)

// UserConsents model is table that stores AA consent details created by the user
type UserConsents struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;PRIMARY_KEY;"`
	UserID            string    `json:"user_id"`
	CustomerID        string    `json:"customer_id"`
	SessionID         string    `json:"session_id"`
	RahasyaNonce      string    `json:"rahasya_nonce"`
	RahasyaPrivateKey string    `json:"rahasya_private_key"`
	Expire            time.Time `json:"expire"`
	ConsentStatus     string    `json:"consent_status"`
	ArtefactStatus    string    `json:"artefact_status"`
	DataFetched       bool      `json:"dataFetched"`
	SignedConsent     string    `json:"signed_consent" gorm:"default:-"`
	ConsentHandle     uuid.UUID `json:"consent_handle" gorm:"type:uuid;"`
	ConsentID         uuid.UUID `json:"consent_id" gorm:"type:uuid;"`
	User              Users     `gorm:"foreignKey:UserID"`
}

// BeforeCreate creates new UUID before saving to database
func (consent *UserConsents) BeforeCreate(tx *gorm.DB) (err error) {
	consent.ID = uuid.New()
	return
}
