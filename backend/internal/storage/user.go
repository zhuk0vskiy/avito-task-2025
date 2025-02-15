package storage

import (
	strgDto "avito-task-2025/backend/internal/storage/dto"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=UserIntf
type UserIntf interface {
	Insert(ctx context.Context, request *strgDto.InsertUserRequest) (response *strgDto.InsertUserResponse, err error)
	GetByUsername(ctx context.Context, request *strgDto.GetUserByUsernameRequest) (response *strgDto.GetUserByUsernameResponse, err error)
	GetCoinsByUserID(ctx context.Context, request *strgDto.GetCoinsByUserIDRequest) (response *strgDto.GetCoinsByUserIDResponse, err error)
}
