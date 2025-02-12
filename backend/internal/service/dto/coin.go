package dto

import "github.com/google/uuid"

type SendCoinsRequest struct {
	UserID         uuid.UUID
	ToUserUsername string
	CoinsAmount    int32
}

// type SendResponse struct {

// }
