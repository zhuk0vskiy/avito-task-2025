package user

import (
	"context"
	svcDto "avito-task-2025/backend/internal/service/dto"
)
type Intf interface {
	GetInfo(ctx context.Context, request *svcDto.GetInfoRequest) (response *svcDto.GetInfoResponse, err error)
}