package v1

import (
	"avito-task-2025/backend/internal/app"
	"avito-task-2025/backend/internal/controller"
	svcDto "avito-task-2025/backend/internal/service/dto"
	"encoding/json"

	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func SendCoinHandler(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//start := time.Now()

		wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

		var httpReq SendCoinRequest

		err := json.NewDecoder(r.Body).Decode(&httpReq)
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
			return
		}

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
		req := &svcDto.SendCoinsRequest{
			UserID:         uuID,
			ToUserUsername: httpReq.ToUser,
			CoinsAmount:    httpReq.Amount,
		}
		err = a.CoinSvcIntf.Send(r.Context(), req)
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, nil)
	}
}
