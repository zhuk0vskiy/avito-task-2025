package v1

type SignInResponse struct {
	Token string `json:"token"`
}

type SendCoinResponse struct {
}

type BuyMerchResponse struct {
}

type Inventory struct {
	Type     string `json:"type"`
	Quantity int16  `json:"quantity"`
}

type CoinHistoryReceived struct {
	FromUser string `json:"fromUser"`
	Amount   int32  `json:"amount"`
}

type CoinHistorySent struct {
	ToUser string `json:"toUser"`
	Amount int32  `json:"amount"`
}

type CoinHistory struct {
	Received []*CoinHistoryReceived `json:"received"`
	Sent     []*CoinHistorySent     `json:"sent"`
}

type GetInfoResponse struct {
	Coins       int32        `json:"coins"`
	Inventory   []*Inventory `json:"inventory"`
	CoinHistory *CoinHistory `json:"coinHistory"`
}

type ErrorResponseStruct struct {
	Error string `json:"errors"`
}
