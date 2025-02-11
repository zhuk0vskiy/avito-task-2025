package storage

import (
	"context"
	strgDto "avito-task-2025/backend/internal/storage/dto"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=TransactionIntf
type TransactionIntf interface {
	Insert(ctx context.Context, request *strgDto.InsertTransactionRequest) (err error)
	GetByFromUserID(ctx context.Context, request *strgDto.GetTransactionByFromUserIDRequest) (response *strgDto.GetTransactionByToUserIDResponse, err error)
	GetByToUserID(ctx context.Context, request *strgDto.GetTransactionByToUserIDRequest) (response *strgDto.GetTransactionByToUserIDResponse, err error)
}