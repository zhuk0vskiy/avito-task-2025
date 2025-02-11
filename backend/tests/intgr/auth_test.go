package intgr

import (
	"context"
	"testing"
	"time"

	loggerMock "avito-task-2025/backend/pkg/logger/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"avito-task-2025/backend/internal/service"
	svcDto "avito-task-2025/backend/internal/service/dto"
	"avito-task-2025/backend/internal/storage/postgres"
	// "avito-task-2025/backend/internal/storage/"
)

// success sign up
func TestSignInSuccess_01(t *testing.T) {

	mockLogger := new(loggerMock.Interface)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg := NewTestConfig()

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	userStrgIntf := postgres.NewUserStrg(dbConnector)

	authService := service.NewAuthSvc(mockLogger, userStrgIntf, cfg.JwtKey)

	req := &svcDto.SignInRequest{
		Username: "test1",
		Password: "test1",
	}

	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Warnf", mock.Anything, mock.Anything).Times(0)

	response, err := authService.SignIn(ctx, req)

	assert.NoError(t, err)
	assert.NotEmpty(t, response)
}


// success sign up and log in
func TestSignInSuccess_02(t *testing.T) {

	mockLogger := new(loggerMock.Interface)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg := NewTestConfig()

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	userStrgIntf := postgres.NewUserStrg(dbConnector)

	authService := service.NewAuthSvc(mockLogger, userStrgIntf, cfg.JwtKey)

	req := &svcDto.SignInRequest{
		Username: "test2",
		Password: "test2",
	}

	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Warnf", mock.Anything, mock.Anything).Times(0)

	_, err = authService.SignIn(ctx, req)
	if err != nil {
		t.Error(err)
	}

	response, err := authService.SignIn(ctx, req)

	assert.NoError(t, err)
	assert.NotEmpty(t, response)
}

// success sign up and failed log in
func TestSignInFailed_01(t *testing.T) {

	mockLogger := new(loggerMock.Interface)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg := NewTestConfig()

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	userStrgIntf := postgres.NewUserStrg(dbConnector)

	authService := service.NewAuthSvc(mockLogger, userStrgIntf, cfg.JwtKey)

	req := &svcDto.SignInRequest{
		Username: "test3",
		Password: "test3",
	}

	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Warnf", mock.Anything, mock.Anything).Times(0)

	_, err = authService.SignIn(ctx, req)
	if err != nil {
		t.Error(err)
	}

	req = &svcDto.SignInRequest{
		Username: "test3",
		Password: "test-",
	}
	response, err := authService.SignIn(ctx, req)

	assert.Error(t, err)
	assert.Empty(t, response)
}



