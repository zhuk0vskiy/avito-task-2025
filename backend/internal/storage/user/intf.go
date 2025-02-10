package user

import (
	"context"
)

type Intf interface{
	Insert(ctx context.Context, request *InsertRequest) (err error)
	GetByUsername(ctx context.Context, request *GetByUsernameRequest) (response *GetByUsernameResponse, err error)
}