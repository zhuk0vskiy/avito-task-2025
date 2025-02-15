package dto

import "github.com/google/uuid"

type SendCoinsRequest struct {
	UserID         uuid.UUID
	ToUsername string
	CoinsAmount    int32
}

// type SendResponse struct {

// }
