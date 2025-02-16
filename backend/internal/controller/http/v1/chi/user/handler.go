package user

import (
	"avito-task-2025/backend/internal/app"
	v1 "avito-task-2025/backend/internal/controller/http/v1"
	svcDto "avito-task-2025/backend/internal/service/dto"
	"encoding/json"

	"net/http"

	"github.com/google/uuid"
)

func GetUserInfoHandler(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//start := time.Now()
		// fmt.Println(r.Context())
		id, err := a.JwtMngIntf.GetStringClaimFromJWT(r.Context(), "id")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(v1.ErrorResponseStruct{Error: "jwt token is invalid"})
			return
		}

		uuID, err := uuid.Parse(id)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(v1.ErrorResponseStruct{Error: err.Error()})
			return
		}
		req := &svcDto.GetUserInfoRequest{
			UserID: uuID,
		}
		svcResponse, err := a.UserSvcIntf.GetInfo(r.Context(), req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(v1.ErrorResponseStruct{Error: err.Error()})
			return
		}

		inventory := make([]*v1.Inventory, 0)
		for i := 0; i < len(svcResponse.Inventory); i++ {
			inventory = append(inventory, &v1.Inventory{
				Type:     svcResponse.Inventory[i].Type,
				Quantity: svcResponse.Inventory[i].Amount,
			})
		}

		receive := make([]*v1.CoinHistoryReceived, 0)
		for i := 0; i < len(svcResponse.CoinHistory.Received); i++ {
			receive = append(receive, &v1.CoinHistoryReceived{
				FromUser: svcResponse.CoinHistory.Received[i].FromUsername,
				Amount:   svcResponse.CoinHistory.Received[i].CoinsAmount,
			})
		}

		sent := make([]*v1.CoinHistorySent, 0)
		for i := 0; i < len(svcResponse.CoinHistory.Sent); i++ {
			sent = append(sent, &v1.CoinHistorySent{
				ToUser: svcResponse.CoinHistory.Sent[i].ToUsername,
				Amount: svcResponse.CoinHistory.Sent[i].CoinsAmount,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&v1.GetInfoResponse{
			Coins:     svcResponse.Coins,
			Inventory: inventory,
			CoinHistory: &v1.CoinHistory{
				Received: receive,
				Sent:     sent,
			},
		})
	}
}
