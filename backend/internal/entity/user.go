package entity

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Username     string
	Password     string
	CoinsAmount  int32
	Transactions []*Transaction
}
