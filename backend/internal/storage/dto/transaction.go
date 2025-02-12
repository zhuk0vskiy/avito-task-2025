package dto

import (
	"avito-task-2025/backend/internal/entity"

	"github.com/google/uuid"
)

type InsertTransactionRequest struct {
	FromUserID  uuid.UUID
	ToUsername  string
	CoinsAmount int32
}

type InsertTransactionResponse struct {
}

type GetTransactionByFromUserIDRequest struct {
	UserID uuid.UUID
}

type GetTransactionByFromUserIDResponse struct {
	Transactions []*entity.Transaction
}

type GetTransactionByToUserIDRequest struct {
	UserID uuid.UUID
}

type GetTransactionByToUserIDResponse struct {
	Transactions []*entity.Transaction
}
