package e2e

import (
	v1 "avito-task-2025/backend/internal/controller/http/v1"

	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func SignInSuccess_01(t *testing.T, baseUrl string) {
	reqBody, err := json.Marshal(v1.SignInRequest{
		Username: "test",
		Password: "test",
	})
	require.NoError(t, err)

	resp, err := http.Post(baseUrl+"/api/signup", "application/json", bytes.NewBuffer(reqBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
