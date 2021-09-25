package insurance

import "github.com/google/uuid"

type purchaseRequest struct {
	Uuid uuid.UUID `json:"uuid"`
}
