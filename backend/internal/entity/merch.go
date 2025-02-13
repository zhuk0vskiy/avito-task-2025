package entity

import "github.com/google/uuid"

type Merch struct {
	ID      uuid.UUID
	OwnerID uuid.UUID
	Type    string
	Cost    int32
	Amount  int16
}
