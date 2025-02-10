package user

type InsertRequest struct {
	Username string
	HashPassword []byte
}

// type InsertResponse struct {

// }

type GetByUsernameRequest struct {
	Username string
}

type GetByUsernameResponse struct {
	HashPassword []byte
}

