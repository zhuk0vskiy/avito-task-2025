package v1

import (
	"avito-task-2025/backend/internal/app"
	"avito-task-2025/backend/internal/controller"
	svcDto "avito-task-2025/backend/internal/service/dto"

	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func GetUserInfoHandler(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//start := time.Now()

		wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

		// var httpReq SignInRequest

		// err := json.NewDecoder(r.Body).Decode(&httpReq)
		// if err != nil {
		// 	controller.ErrorResponse(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		// 	return
		// }

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
