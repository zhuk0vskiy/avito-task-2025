package auth

import (
	"avito-task-2025/backend/internal/app"
	"avito-task-2025/backend/internal/controller"
	svcDto "avito-task-2025/backend/internal/service/dto"

	"encoding/json"
	"fmt"
	"net/http"
)

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

func SignInHandler(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//start := time.Now()

		wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

		var httpReq SignInRequest

		err := json.NewDecoder(r.Body).Decode(&httpReq)
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
			return
		}

		req := &svcDto.SignInRequest{
			Username: httpReq.Username,
			Password: httpReq.Password,
		}
		svcResponse, err := a.AuthSvcIntf.SignIn(r.Context(), req)
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
			return
		}

		response := &SignInResponse{
			Token: svcResponse.JwtToken,
		}
		controller.SuccessResponse(wrappedWriter, http.StatusOK, map[string]string{"token": response.Token})
	}
}
