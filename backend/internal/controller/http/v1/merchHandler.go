package v1

import (
	"avito-task-2025/backend/internal/app"
	"avito-task-2025/backend/internal/controller"
	svcDto "avito-task-2025/backend/internal/service/dto"

	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func BuyMerchHandler(a *app.App) http.HandlerFunc {
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

		item := chi.URLParam(r, "item")

		req := &svcDto.BuyMerchRequest{
			UserID:    uuID,
			MerchName: item,
		}
		err = a.MerchSvcIntf.Buy(r.Context(), req)
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, nil)
	}
}
