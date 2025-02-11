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

// success send
func TestSendCoinsSuccess_01(t *testing.T) {

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

	authSvcIntf := service.NewAuthSvc(mockLogger, userStrgIntf, cfg.JwtKey)
	coinSvcIntf := service.NewCoinSvc(mockLogger, transactionStrgIntf)

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
	
	err = coinSvcIntf.Send(ctx, &svcDto.SendCoinsRequest{
		UserID: id,
		ToUserUsername: "test",
		CoinsAmount: 10,
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, response)
}

// incorrect user id
func TestSendCoinsFailed_01(t *testing.T) {

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

	authSvcIntf := service.NewAuthSvc(mockLogger, userStrgIntf, cfg.JwtKey)
	coinSvcIntf := service.NewCoinSvc(mockLogger, transactionStrgIntf)

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

	_, err = jwt.VerifyAuthToken(token, cfg.JwtKey)
	if err != nil {
		t.Error(err)
	}
	id, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}
	
	err = coinSvcIntf.Send(ctx, &svcDto.SendCoinsRequest{
		UserID: id,
		ToUserUsername: "test",
		CoinsAmount: 10,
	})

	assert.Error(t, err)
}


// incorrect receiver
func TestSendCoinsFailed_02(t *testing.T) {

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

	authSvcIntf := service.NewAuthSvc(mockLogger, userStrgIntf, cfg.JwtKey)
	coinSvcIntf := service.NewCoinSvc(mockLogger, transactionStrgIntf)

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
	
	err = coinSvcIntf.Send(ctx, &svcDto.SendCoinsRequest{
		UserID: id,
		ToUserUsername: "avito_test",
		CoinsAmount: 10,
	})
	
	assert.Error(t, err)
}

// coins amount less then 0
func TestSendCoinsFailed_03(t *testing.T) {

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

	authSvcIntf := service.NewAuthSvc(mockLogger, userStrgIntf, cfg.JwtKey)
	coinSvcIntf := service.NewCoinSvc(mockLogger, transactionStrgIntf)

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
	
	err = coinSvcIntf.Send(ctx, &svcDto.SendCoinsRequest{
		UserID: id,
		ToUserUsername: "avito_test",
		CoinsAmount: -1,
	})
	
	assert.Error(t, err)
}





