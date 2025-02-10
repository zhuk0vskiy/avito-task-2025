package boughtmerch

import "github.com/google/uuid"

type AddRequest struct {
	UserID    uuid.UUID
	MerchName string
}

// type AddResponse struct {

// }

type GetByUserIDRequest struct {
	UserID uuid.UUID
}

type GetByUserIDResponse struct {
	Merchs []*struct {
		Name   string
		Amount int16
	}
}
