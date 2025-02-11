package dto

import (
	"avito-task-2025/backend/internal/entity"

	"github.com/google/uuid"
)

type GetUserInfoRequest struct {
	UserID uuid.UUID
}

type GetUserInfoResponse struct {
	UserInfo *struct {
		Coins       int32
		Inventory   []*entity.Merch
		CoinHistory *struct {
			Received []*entity.Transaction
			Sent     []*entity.Transaction
		}
	}
}
