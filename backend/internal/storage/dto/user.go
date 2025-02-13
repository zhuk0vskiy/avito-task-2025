package dto

import "github.com/google/uuid"

type InsertUserRequest struct {
	Username     string
	HashPassword []byte
}

// type InsertUserResponse struct {

// }

type GetUserByUsernameRequest struct {
	Username string
}

type GetUserByUsernameResponse struct {
	ID           uuid.UUID
	HashPassword []byte
}

type GetCoinsByUserIDRequest struct {
	UserID uuid.UUID
}

type GetCoinsByUserIDResponse struct {
	Amount int32
}
