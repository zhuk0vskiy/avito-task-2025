package dto

import "github.com/google/uuid"

type SendRequest struct {
	ToUserID uuid.UUID
	CoinsAmount int32
}

// type SendResponse struct {

// }