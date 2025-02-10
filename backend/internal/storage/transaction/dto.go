package transaction

import (
	"avito-task-2025/backend/internal/entity"

	"github.com/google/uuid"
)

type InsertRequest struct {
	FromUserID  uuid.UUID
	ToUserID    string
	CoinsAmount int32
}

type InsertResponse struct {

}

type GetByFromUserIDRequest struct {
	UserID uuid.UUID
}

type GetByFromUserIDResponse struct {
	Transactions []*entity.Transaction
}

type GetByToUserIDRequest struct {
	UserID uuid.UUID
}

type GetByToUserIDResponse struct {
	Transactions []*entity.Transaction
}
