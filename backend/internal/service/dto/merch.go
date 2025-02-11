package dto

import "github.com/google/uuid"

type BuyMerchRequest struct {
	UserID uuid.UUID
	MerchName string
}

// type BuyResponse struct {

// }