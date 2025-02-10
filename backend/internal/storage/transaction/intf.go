package transaction

import "context"

type Intf interface {
	Insert(ctx context.Context, request *InsertRequest) (err error)
	GetByFromUserID(ctx context.Context, request *GetByFromUserIDRequest) (response *GetByFromUserIDResponse, err error)
	GetByToUserID(ctx context.Context, request *GetByToUserIDRequest) (response *GetByToUserIDResponse, err error)
}