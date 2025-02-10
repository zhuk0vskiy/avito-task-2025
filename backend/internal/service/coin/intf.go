package coin


import (
	svcDto "avito-task-2025/backend/internal/service/dto"
	"context"
)

type Intf interface {
	Send(ctx context.Context, request *svcDto.BuyRequest) (err error)
}