package coin

import (
	"avito-task-2025/backend/internal/app"

	v1 "avito-task-2025/backend/internal/controller/http/v1"
	svcDto "avito-task-2025/backend/internal/service/dto"
	"encoding/json"

	"net/http"

	"github.com/google/uuid"
)

func SendCoinHandler(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var httpReq v1.SendCoinRequest

		err := json.NewDecoder(r.Body).Decode(&httpReq)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(v1.ErrorResponseStruct{Error: err.Error()})
			return
		}

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
		req := &svcDto.SendCoinsRequest{
			UserID:      uuID,
			ToUsername:  httpReq.ToUser,
			CoinsAmount: httpReq.Amount,
		}
		err = a.CoinSvcIntf.Send(r.Context(), req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(v1.ErrorResponseStruct{Error: err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
