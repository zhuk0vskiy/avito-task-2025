package dto

import (
	"avito-task-2025/backend/internal/entity"

	"github.com/google/uuid"
)

type InsertBoughtMerchRequest struct {
	UserID uuid.UUID
	Type   string
}

// type AddBoughtMerchResponse struct {

// }

type GetBoughtMerchByUserIDRequest struct {
	UserID uuid.UUID
}

type GetBoughtMerchByUserIDResponse struct {
	Merchs []*entity.Merch
}
