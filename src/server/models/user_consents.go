package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	EmptyColumn       = "-"
	StatusDataFetched = "FETCHED"
)

const (
	ArtefactStatusReady   = "READY"
	ArtefactStatusActive  = "ACTIVE"
	ArtefactStatusDenied  = "DENIED"
	ArtefactStatusPending = "PENDING"
	ArtefactStatusTimeout = "TIMEOUT"
	ArtefactStatusError   = "UNKNOWN"
)

const (
	ConsentStatusActive   = "ACTIVE"
	ConsentStatusRejected = "REJECTED"
	ConsentStatusRevoked  = "REVOKED"
	ConsentStatusPaused   = "PAUSED"
	ConsentStatusError    = "UNKNOWN"
)

type UserConsents struct {
	Id                uuid.UUID `json:"id" gorm:"type:uuid;PRIMARY_KEY;"`
	UserId            string    `json:"user_id"`
	CustomerId        string    `json:"customer_id"`
	SessionId         string    `json:"session_id"`
	RahasyaNonce      string    `json:"rahasya_nonce"`
	RahasyaPrivateKey string    `json:"rahasya_private_key"`
	Expire            time.Time `json:"expire"`
	ConsentStatus     string    `json:"consent_status"`
	ArtefactStatus    string    `json:"artefact_status"`
	DataFetched       bool      `json:"dataFetched"`
	SignedConsent     string    `json:"signed_consent" gorm:"default:-"`
	ConsentHandle     uuid.UUID `json:"consent_handle" gorm:"type:uuid;"`
	ConsentId         uuid.UUID `json:"consent_id" gorm:"type:uuid;"`
	User              Users     `gorm:"foreignKey:UserId"`
}

func (consent *UserConsents) BeforeCreate(tx *gorm.DB) (err error) {
	consent.Id = uuid.New()
	return
}
