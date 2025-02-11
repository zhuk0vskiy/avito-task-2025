package intgr

import (
	"avito-task-2025/backend/internal/service"
	svcDto "avito-task-2025/backend/internal/service/dto"
	"avito-task-2025/backend/internal/storage/postgres"
	"avito-task-2025/backend/pkg/jwt"
	loggerMock "avito-task-2025/backend/pkg/logger/mocks"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// success to buy merch
func TestMerchBuySuccess_01(t *testing.T) {

	mockLogger := new(loggerMock.Interface)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg := NewTestConfig()

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	userStrgIntf := postgres.NewUserStrg(dbConnector)
	boughtMerchStrgIntf := postgres.NewBoughtMerchStrg(dbConnector)

	authSvcIntf := service.NewAuthSvc(mockLogger, userStrgIntf, cfg.JwtKey)
	merchSvcIntf := service.NewMerchSvc(mockLogger, boughtMerchStrgIntf)

	req := &svcDto.SignInRequest{
		Username: "test5",
		Password: "test5",
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

	err = merchSvcIntf.Buy(ctx, &svcDto.BuyMerchRequest{
		UserID:    id,
		MerchName: "t-shirt",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, response)
}

// incorrect user id
func TestMerchBuyFailed_01(t *testing.T) {

	mockLogger := new(loggerMock.Interface)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg := NewTestConfig()

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	userStrgIntf := postgres.NewUserStrg(dbConnector)
	boughtMerchStrgIntf := postgres.NewBoughtMerchStrg(dbConnector)

	authSvcIntf := service.NewAuthSvc(mockLogger, userStrgIntf, cfg.JwtKey)
	merchSvcIntf := service.NewMerchSvc(mockLogger, boughtMerchStrgIntf)

	req := &svcDto.SignInRequest{
		Username: "test5",
		Password: "test5",
	}

	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Warnf", mock.Anything, mock.Anything).Times(0)

	response, err := authSvcIntf.SignIn(ctx, req)
	if err != nil {
		t.Error(err)
	}
	token := response.JwtToken

	_, err = jwt.VerifyAuthToken(token, cfg.JwtKey)
	if err != nil {
		t.Error(err)
	}
	id, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}

	err = merchSvcIntf.Buy(ctx, &svcDto.BuyMerchRequest{
		UserID:    id,
		MerchName: "t-shirt",
	})

	assert.Error(t, err)
	assert.Equal(t, postgres.Test_errBuyDecreaseCoinsAmount, err)
}

// incorrect merch type
func TestMerchBuyFailed_02(t *testing.T) {

	mockLogger := new(loggerMock.Interface)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg := NewTestConfig()

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	userStrgIntf := postgres.NewUserStrg(dbConnector)
	boughtMerchStrgIntf := postgres.NewBoughtMerchStrg(dbConnector)

	authSvcIntf := service.NewAuthSvc(mockLogger, userStrgIntf, cfg.JwtKey)
	merchSvcIntf := service.NewMerchSvc(mockLogger, boughtMerchStrgIntf)

	req := &svcDto.SignInRequest{
		Username: "test5",
		Password: "test5",
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

	err = merchSvcIntf.Buy(ctx, &svcDto.BuyMerchRequest{
		UserID:    id,
		MerchName: "another avito merch",
	})

	assert.Error(t, err)
	assert.Equal(t, postgres.Test_errBuyDecreaseCoinsAmount, err)
}
