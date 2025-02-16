package storage

import (
	strgDto "avito-task-2025/backend/internal/storage/dto"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=BoughtMerchIntf
type BoughtMerchIntf interface {
	Insert(ctx context.Context, request *strgDto.InsertBoughtMerchRequest) (err error)
	GetByUserID(ctx context.Context, request *strgDto.GetBoughtMerchByUserIDRequest) (response *strgDto.GetBoughtMerchByUserIDResponse, err error)
}
