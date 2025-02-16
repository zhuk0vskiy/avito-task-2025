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

type GetTransactionToUserIDRequest struct {
	UserID uuid.UUID
}

type GetTransactionToUserIDResponse struct {
	Transactions []*entity.Transaction
}

type GetTransactionFromUserIDRequest struct {
	UserID uuid.UUID
}

type GetTransactionFromUserIDResponse struct {
	Transactions []*entity.Transaction
}
