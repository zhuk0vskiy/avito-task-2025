package auth

import (
	"avito-task-2025/backend/internal/app"
	v1 "avito-task-2025/backend/internal/controller/http/v1"

	// "avito-task-2025/backend/internal/controller"
	svcDto "avito-task-2025/backend/internal/service/dto"

	"encoding/json"
	"net/http"
)

func SignInHandler(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var httpReq v1.SignInRequest

		err := json.NewDecoder(r.Body).Decode(&httpReq)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(v1.ErrorResponseStruct{Error: err.Error()})
			return
		}

		req := &svcDto.SignInRequest{
			Username: httpReq.Username,
			Password: httpReq.Password,
		}
		svcResponse, err := a.AuthSvcIntf.SignIn(r.Context(), req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(v1.ErrorResponseStruct{Error: err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&v1.SignInResponse{
			Token: svcResponse.JwtToken,
		})
	}
}
