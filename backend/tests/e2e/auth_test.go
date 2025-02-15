package e2e

import (
	"avito-task-2025/backend/internal/controller/http/v1"
	"avito-task-2025/backend/tests"
	"net/http"
	"testing"
)

func TestAuthRequests(t *testing.T) {
	migrator, _ := tests.NewTestConfig("file://../../../db/postgres/test_migrations/e2e/")

	_ = migrator.Force(1)
	_ = migrator.Down()

	e := tests.NewExpect(t)

	req := &v1.SignInRequest{
		Username: "test",
		Password: "test",
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

	failReq := &v1.SignInRequest{
		Username: "test",
		Password: "test1",
	}

	_ = e.POST("/auth").
		WithJSON(failReq).
		Expect().
		Status(http.StatusBadRequest).
		JSON().
		Object().
		NotEmpty().
		ContainsKey("errors").
		HasValue("errors", "incorrect password")

	failReq1 := &v1.SignInRequest{
		Username: "",
		Password: "test1",
	}

	_ = e.POST("/auth").
		WithJSON(failReq1).
		Expect().
		Status(http.StatusBadRequest).
		JSON().
		Object().
		NotEmpty().
		ContainsKey("errors").
		HasValue("errors", "username is empty")

	failReq2 := &v1.SignInRequest{
		Username: "r",
		Password: "",
	}

	_ = e.POST("/auth").
		WithJSON(failReq2).
		Expect().
		Status(http.StatusBadRequest).
		JSON().
		Object().
		NotEmpty().
		ContainsKey("errors").
		HasValue("errors", "password is empty")
}
