package user

import (
	"avito-task-2025/backend/internal/app"
	"avito-task-2025/backend/internal/controller"
	svcDto "avito-task-2025/backend/internal/service/dto"

	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type GetInfoRequest struct {

}

type GetInfoInventory struct {
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
	Coins       int32               `json:"coins"`
	Inventory   []*GetInfoInventory `json:"inventory"`
	CoinHistory *CoinHistory        `json:"coinHistory"`
}

func GetUserInfoHandler(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//start := time.Now()

		wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}


		id, err := a.JwtIntf.GetStringClaimFromJWT(r.Context(), "id")
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s", err).Error(), http.StatusInternalServerError)
			return
		}

		uuID, err := uuid.Parse(id)
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s", err).Error(), http.StatusInternalServerError)
			return
		}
		req := &svcDto.GetUserInfoRequest{
			UserID: uuID,
		}
		svcResponse, err := a.UserSvcIntf.GetInfo(r.Context(), req)
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
			return
		}

		// response := &GetInfoResponse{
		// 	Coins:     svcResponse.Coins,
		// 	Inventory: svcResponse.Inventory,
		// 	CoinHistory: &struct {
		// 		Received []*entity.Transaction
		// 		Sent     []*entity.Transaction
		// 	}{
		// 		Received: svcResponse.CoinHistory.Received,
		// 		Sent:     svcResponse.CoinHistory.Sent,
		// 	},
		// }

		inventory := make([]*GetInfoInventory, 0)
		for i := 0; i < len(svcResponse.Inventory); i++ {
			inventory = append(inventory, &GetInfoInventory{
				Type:     svcResponse.Inventory[i].Type,
				Quantity: svcResponse.Inventory[i].Amount,
			})
		}

		receive := make([]*CoinHistoryReceived, 0)
		for i := 0; i < len(svcResponse.CoinHistory.Received); i++ {
			receive = append(receive, &CoinHistoryReceived{
				FromUser: svcResponse.CoinHistory.Received[i].FromUsername,
				Amount:   svcResponse.CoinHistory.Received[i].CoinsAmount,
			})
		}

		sent := make([]*CoinHistorySent, 0)
		for i := 0; i < len(svcResponse.CoinHistory.Sent); i++ {
			sent = append(sent, &CoinHistorySent{
				ToUser: svcResponse.CoinHistory.Sent[i].ToUsername,
				Amount: svcResponse.CoinHistory.Sent[i].CoinsAmount,
			})
		}

		response := &GetInfoResponse{
			Coins:     svcResponse.Coins,
			Inventory: inventory,
			CoinHistory: &CoinHistory{
				Received: receive,
				Sent:     sent,
			},
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, response)
	}
}
