package entity

import "github.com/google/uuid"

type Transaction struct {
	ID          uuid.UUID
	FromUsername  string
	ToUsername    string
	CoinsAmount int32
}
