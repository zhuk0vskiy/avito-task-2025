package e2e

import (
	"avito-task-2025/backend/internal/controller/http/v1"
	"avito-task-2025/backend/tests"
	"net/http"
	"testing"
)

func TestUserRequests(t *testing.T) {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	migrator, _ := tests.NewTestConfig("file://../../../db/postgres/test_migrations/e2e/")

	_ = migrator.Force(1)
	_ = migrator.Down()

	e := tests.NewExpect(t)

	req := &v1.SignInRequest{
		Username: "test1",
		Password: "test1",
	}
	token := e.POST("/auth").
		WithJSON(req).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		NotEmpty().
		ContainsKey("token").
		Value("token").Raw().(string)

	req = &v1.SignInRequest{
		Username: "test2",
		Password: "test2",
	}
	_ = e.POST("/auth").
		WithJSON(req).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		NotEmpty().
		ContainsKey("token").
		Value("token").Raw().(string)

	sendReq := v1.SendCoinRequest{
		ToUser: "test2",
		Amount: 501,
	}
	_ = e.POST("/sendCoin").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(sendReq).
		Expect().
		Status(http.StatusOK)

	sendReq = v1.SendCoinRequest{
		ToUser: "test3",
		Amount: 501,
	}
	_ = e.POST("/sendCoin").
		WithHeader("Authorization", "Bearer "+"a1").
		WithJSON(sendReq).
		Expect().
		Status(http.StatusBadRequest).
		JSON().
		Object().
		NotEmpty().
		ContainsKey("errors").
		HasValue("errors", "jwt token is invalid")

	sendReq = v1.SendCoinRequest{
		ToUser: "test3",
		Amount: 501,
	}
	_ = e.POST("/sendCoin").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(sendReq).
		Expect().
		Status(http.StatusBadRequest).
		JSON().
		Object().
		NotEmpty().
		ContainsKey("errors").
		HasValue("errors", "receiver with this username doesnt exist")

	sendReq = v1.SendCoinRequest{
		ToUser: "test2",
		Amount: 1000,
	}
	_ = e.POST("/sendCoin").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(sendReq).
		Expect().
		Status(http.StatusBadRequest).
		JSON().
		Object().
		NotEmpty().
		ContainsKey("errors").
		HasValue("errors", "not enough coins to send")
}
