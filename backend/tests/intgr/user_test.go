package intgr

import (
	"avito-task-2025/backend/internal/service"
	svcDto "avito-task-2025/backend/internal/service/dto"
	"avito-task-2025/backend/internal/storage/postgres"
	"avito-task-2025/backend/pkg/jwt"
	loggerMock "avito-task-2025/backend/pkg/logger/mocks"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// success get user
func TestGetUserSuccess_01(t *testing.T) {

	mockLogger := new(loggerMock.Interface)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg := NewTestConfig()

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	userStrgIntf := postgres.NewUserStrg(dbConnector)

	transactionStrgIntf := postgres.NewTransactionStrg(dbConnector)
	boughtMerchStrgIntf := postgres.NewBoughtMerchStrg(dbConnector)



	authSvcIntf := service.NewAuthSvc(mockLogger, userStrgIntf, cfg.JwtKey)
	userSvcIntf := service.NewUserSvc(mockLogger, userStrgIntf, boughtMerchStrgIntf, transactionStrgIntf)

	req := &svcDto.SignInRequest{
		Username: "test6",
		Password: "test6",
	}

	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Warnf", mock.Anything, mock.Anything).Times(0)

	response, err := authSvcIntf.SignIn(ctx, req)
	if err != nil {
		t.Error(err)
	}
	token := response.JwtToken

	payload, err := jwt.VerifyAuthToken(token, cfg.JwtKey)
	if err != nil {
		t.Error(err)
	}
	id, err := uuid.Parse(payload.Id)
	if err != nil {
		t.Error(err)
	}
	
	getResponse, err := userSvcIntf.GetInfo(ctx, &svcDto.GetUserInfoRequest{
		UserID: id,
	})
	if err != nil {
		t.Error(err)
	}

	fmt.Println(getResponse.UserInfo.Coins)
	for _, i := range getResponse.UserInfo.Inventory {
		fmt.Println(i)
	}
	for _, i := range getResponse.UserInfo.CoinHistory.Received {
		fmt.Println(i.FromUsername)
	}
	for _, i := range getResponse.UserInfo.CoinHistory.Sent {
		fmt.Println(i.ToUsername)
	}
	

	assert.NoError(t, err)
	assert.NotEmpty(t, getResponse)
}