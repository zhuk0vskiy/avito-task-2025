package e2e

import (
	"avito-task-2025/backend/internal/controller/http/v1"
	"avito-task-2025/backend/tests"
	"net/http"
	"testing"
)

func TestCoinRequests(t *testing.T) {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	migrator, _ := tests.NewTestConfig("file://../../../db/postgres/test_migrations/e2e/")

	_ = migrator.Force(1)
	_ = migrator.Down()

	e := tests.NewExpect(t)

	req := &v1.SignInRequest{
		Username: "test",
		Password: "test",
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

	_ = e.GET("/info").
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		NotEmpty().
		ContainsKey("coins").
		ContainsKey("inventory").
		ContainsKey("coinHistory").
		HasValue("coins", 1000)

	_ = e.GET("/info").
		WithHeader("Authorization", "Bearer "+"ac5").
		Expect().
		Status(http.StatusBadRequest).
		JSON().
		Object().
		NotEmpty().
		ContainsKey("errors").
		HasValue("errors", "jwt token is invalid")

	_ = e.GET("/buy/pen").
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK)

	req = &v1.SignInRequest{
		Username: "test2",
		Password: "test2",
	}
	token2 := e.POST("/auth").
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
		Amount: 100,
	}
	_ = e.POST("/sendCoin").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(sendReq).
		Expect().
		Status(http.StatusOK)

	sendReq = v1.SendCoinRequest{
		ToUser: "test",
		Amount: 100,
	}
	_ = e.POST("/sendCoin").
		WithHeader("Authorization", "Bearer "+token2).
		WithJSON(sendReq).
		Expect().
		Status(http.StatusOK)

	inventory := make([]*v1.Inventory, 0)
	inventory = append(inventory, &v1.Inventory{
		Type:     "pen",
		Quantity: 1,
	})
	coinHistoryReceived := make([]*v1.CoinHistoryReceived, 0)
	coinHistoryReceived = append(coinHistoryReceived, &v1.CoinHistoryReceived{
		FromUser: "test2",
		Amount:   100,
	})
	coinHistorySent := make([]*v1.CoinHistorySent, 0)
	coinHistorySent = append(coinHistorySent, &v1.CoinHistorySent{
		ToUser: "test2",
		Amount: 100,
	})

	coinHistory := v1.CoinHistory{
		Received: coinHistoryReceived,
		Sent:     coinHistorySent,
	}

	_ = e.GET("/info").
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		NotEmpty().
		ContainsKey("coins").
		ContainsKey("inventory").
		ContainsKey("coinHistory").
		HasValue("coins", 990).
		HasValue("inventory", inventory).
		HasValue("coinHistory", coinHistory)

	// _ = e.GET("/buy/pen").
	// 	WithHeader("Authorization", "Bearer "+token).
	// 	Expect().
	// 	Status(http.StatusOK)

}
