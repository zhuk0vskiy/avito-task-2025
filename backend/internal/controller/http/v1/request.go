package v1

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int32  `json:"amount"`
}

type BuyMerchRequest struct {
	Item string `json:"item"`
}

type GetInfoRequest struct {
}