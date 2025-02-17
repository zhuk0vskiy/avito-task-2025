package merch

import (
	"avito-task-2025/backend/internal/app"
	v1 "avito-task-2025/backend/internal/controller/http/v1"
	svcDto "avito-task-2025/backend/internal/service/dto"
	"encoding/json"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func BuyMerchHandler(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//start := time.Now()

		item := chi.URLParam(r, "item")

		id, err := a.JwtMngIntf.GetStringClaimFromJWT(r.Context(), "id")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
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

		req := &svcDto.BuyMerchRequest{
			UserID:    uuID,
			MerchName: item,
		}

		err = a.MerchSvcIntf.Buy(r.Context(), req)
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
