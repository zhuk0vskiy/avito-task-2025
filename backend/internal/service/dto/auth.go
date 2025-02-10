package dto

type SignInRequest struct {
	Username string
	Password string
}

type SignInResponse struct {
	JwtToken string
}
