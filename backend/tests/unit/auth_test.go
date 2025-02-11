package unit

import (
	"context"
	"time"

	"avito-task-2025/backend/internal/service"
	svcDto "avito-task-2025/backend/internal/service/dto"
	strgDto "avito-task-2025/backend/internal/storage/dto"
	"avito-task-2025/backend/internal/storage/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	loggerMock "avito-task-2025/backend/pkg/logger/mocks"

	"testing"
)

// create new user
func TestSignInSuccess_01(t *testing.T) {

	mockUserStrgIntf := new(mocks.UserIntf)
	mockLogger := new(loggerMock.Interface)

	authService := service.NewAuthSvc(mockLogger, mockUserStrgIntf, "loveavito")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &svcDto.SignInRequest{
		Username: "1",
		Password: "1",
	}

	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Warnf", mock.Anything, mock.Anything).Times(0)

	id, _ := uuid.NewRandom()
	mockUserStrgIntf.On("GetByUsername", ctx, mock.Anything).Return(&strgDto.GetUserByUsernameResponse{
		ID:           id,
		HashPassword: []byte{},
	}, nil).Run(
		func(args mock.Arguments) {
			r := args.Get(1).(*strgDto.GetUserByUsernameRequest)
			assert.Equal(t, req.Username, r.Username)
		})

	mockUserStrgIntf.On("Insert", ctx, mock.Anything).Return(nil).Run(
		func(args mock.Arguments) {
			r := args.Get(1).(*strgDto.InsertUserRequest)
			assert.Equal(t, req.Username, r.Username)
			assert.NotEmpty(t, r.HashPassword)
			assert.NotEmpty(t, r.CoinsAmount)
		})

	response, err := authService.SignIn(ctx, req)

	assert.NoError(t, err)
	assert.NotEmpty(t, response)
}

// user exist, but password incorrect
func TestSignInFailed_01(t *testing.T) {

	mockUserStrgIntf := new(mocks.UserIntf)
	mockLogger := new(loggerMock.Interface)

	authService := service.NewAuthSvc(mockLogger, mockUserStrgIntf, "loveavito")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &svcDto.SignInRequest{
		Username: "1",
		Password: "1",
	}

	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Warnf", mock.Anything, mock.Anything).Times(0)

	id, _ := uuid.NewRandom()
	mockUserStrgIntf.On("GetByUsername", ctx, mock.Anything).Return(&strgDto.GetUserByUsernameResponse{
		ID:           id,
		HashPassword: []byte{'1'},
	}, nil).Run(
		func(args mock.Arguments) {
			r := args.Get(1).(*strgDto.GetUserByUsernameRequest)
			assert.Equal(t, req.Username, r.Username)
		})

	mockUserStrgIntf.On("Insert", ctx, mock.Anything).Return(nil).Run(
		func(args mock.Arguments) {
			r := args.Get(1).(*strgDto.InsertUserRequest)
			assert.Equal(t, req.Username, r.Username)
			assert.NotEmpty(t, r.HashPassword)
			assert.NotEmpty(t, r.CoinsAmount)
		})

	response, err := authService.SignIn(ctx, req)

	assert.Error(t, err)
	assert.Empty(t, response)
}

// empty username
func TestSignInFailed_02(t *testing.T) {

	mockUserStrgIntf := new(mocks.UserIntf)
	mockLogger := new(loggerMock.Interface)

	authService := service.NewAuthSvc(mockLogger, mockUserStrgIntf, "lovaavito")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &svcDto.SignInRequest{
		Username: "",
		Password: "1",
	}

	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Warnf", mock.Anything, mock.Anything).Times(0)

	response, err := authService.SignIn(ctx, req)

	assert.Error(t, err)
	assert.Empty(t, response)
}

//empty password
func TestSignInFailed_03(t *testing.T) {

	mockUserStrgIntf := new(mocks.UserIntf)
	mockLogger := new(loggerMock.Interface)

	authService := service.NewAuthSvc(mockLogger, mockUserStrgIntf, "loveavito")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &svcDto.SignInRequest{
		Username: "1",
		Password: "",
	}

	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Warnf", mock.Anything, mock.Anything).Times(0)

	response, err := authService.SignIn(ctx, req)

	assert.Error(t, err)
	assert.Empty(t, response)
}
