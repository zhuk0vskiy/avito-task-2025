package merch

import (
	"context"
	svcDto "avito-task-2025/backend/internal/service/dto"
)

type Intf interface{
	Buy(ctx context.Context, request *svcDto.BuyRequest) (err error)
}