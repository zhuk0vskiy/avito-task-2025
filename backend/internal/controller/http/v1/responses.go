package v1

type SignInResponse struct {
	Token string `json:"token"`
}

type BuyItemResponse struct {
}

type SendCoinResponse struct {
}

type GetInfoInventory struct {
	Type     string `json:"type"`
	Quantity uint16 `json:"quantity"`
}

type CoinHistoryReceived struct {
	FromUser string `json:"fromUser"`
	Amount   uint32 `json:"amount"`
}

type CoinHistorySent struct {
	ToUser string `json:"toUser"`
	Amount uint32 `json:"amount"`
}

type CoinHistory struct {
	Received []*CoinHistoryReceived `json:"received"`
	Sent     []*CoinHistorySent     `json:"sent"`
}

type GetInfoResponse struct {
	Coins       uint32              `json:"coins"`
	Inventory   []*GetInfoInventory `json:"inventory"`
	CoinHistory *CoinHistory        `json:"coinHistory"`
}
