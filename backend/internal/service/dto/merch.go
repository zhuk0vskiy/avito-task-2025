package dto

import "github.com/google/uuid"

type BuyRequest struct {
	Username uuid.UUID
	MerchName string
}

// type BuyResponse struct {

// }