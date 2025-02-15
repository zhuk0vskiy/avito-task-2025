package storage

import (
	strgDto "avito-task-2025/backend/internal/storage/dto"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=TransactionIntf
type TransactionIntf interface {
	Insert(ctx context.Context, request *strgDto.InsertTransactionRequest) (err error)
	GetToUserID(ctx context.Context, request *strgDto.GetTransactionToUserIDRequest) (response *strgDto.GetTransactionToUserIDResponse, err error)
	GetFromUserID(ctx context.Context, request *strgDto.GetTransactionFromUserIDRequest) (response *strgDto.GetTransactionFromUserIDResponse, err error)
}
