package auth

import (
	svcDto "avito-task-2025/backend/internal/service/dto"
	"context"
)

type Intf interface {
	SignIn(ctx context.Context, request *svcDto.SignInRequest) (response *svcDto.SignInResponse, err error)
}
