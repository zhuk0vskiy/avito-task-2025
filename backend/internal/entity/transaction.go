package entity

import "github.com/google/uuid"

type Transaction struct {
	ID          uuid.UUID
	FromUserID  uuid.UUID
	ToUserID    uuid.UUID
	CoinsAmount int32
}
