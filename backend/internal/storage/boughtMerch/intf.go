package boughtmerch

import (
	"context"
)

type Intf interface {
	Insert(ctx context.Context, request *AddRequest) (err error)
	GetByUserID(ctx context.Context, request *GetByUserIDRequest) (response *GetByUserIDResponse, err error)

}