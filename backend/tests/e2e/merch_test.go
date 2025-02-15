package e2e

import (
	"avito-task-2025/backend/internal/controller/http/v1"
	"avito-task-2025/backend/tests"
	"net/http"
	"testing"
)

func TestMerchRequests(t *testing.T) {

	migrator, _ := tests.NewTestConfig("file://../../../db/postgres/test_migrations/e2e/")

	_ = migrator.Force(1)
	_ = migrator.Down()

	e := tests.NewExpect(t)

	req := &v1.SignInRequest{
		Username: "test2",
		Password: "test2",
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

	_ = e.GET("/buy/pen").
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK)

	_ = e.GET("/buy/pen").
		WithHeader("Authorization", "Bearer "+"a1").
		Expect().
		Status(http.StatusBadRequest).
		JSON().
		Object().
		NotEmpty().
		ContainsKey("errors").
		HasValue("errors", "jwt token is invalid")

	_ = e.GET("/buy/p").
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusBadRequest).
		JSON().
		Object().
		NotEmpty().
		ContainsKey("errors").
		HasValue("errors", "this merch doesnt exist")

	_ = e.GET("/buy/pink-hoody").
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK)

	_ = e.GET("/buy/pink-hoody").
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusBadRequest).
		JSON().
		Object().
		NotEmpty().
		ContainsKey("errors").
		HasValue("errors", "not enough coins to buy merch")
}
