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

	"avito-task-2025/backend/pkg/jwt"
	loggerMock "avito-task-2025/backend/pkg/logger/mocks"

	"testing"
)

func TestSignUpUser(t *testing.T) {

	mockUserStrgIntf := new(mocks.UserIntf)
	mockLogger := new(loggerMock.Interface)
	jwtMngIntf := jwt.NewJwtManager("avito", 1)

	authService := service.NewAuthSvc(mockLogger, jwtMngIntf, mockUserStrgIntf)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &svcDto.SignInRequest{
		Username: "1",
		Password: "1",
	}

	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Warnf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Times(0)

	id, _ := uuid.NewRandom()
	mockUserStrgIntf.On("GetByUsername", ctx, mock.Anything).Return(&strgDto.GetUserByUsernameResponse{
		ID:           id,
		HashPassword: []byte{},
	}, nil)

	mockUserStrgIntf.On("Insert", ctx, mock.Anything).Return(&strgDto.InsertUserResponse{
		ID: id,
	}, nil)

	response, err := authService.SignIn(ctx, req)

	assert.NoError(t, err)
	assert.NotEmpty(t, response)
}

// user exist, but password incorrect
func TestSignInFailed_01(t *testing.T) {

	mockUserStrgIntf := new(mocks.UserIntf)
	mockLogger := new(loggerMock.Interface)
	jwtMngIntf := jwt.NewJwtManager("avito", 1)

	authService := service.NewAuthSvc(mockLogger, jwtMngIntf, mockUserStrgIntf)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &svcDto.SignInRequest{
		Username: "1",
		Password: "1",
	}

	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Warnf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Times(0)

	id, _ := uuid.NewRandom()
	mockUserStrgIntf.On("GetByUsername", ctx, mock.Anything).Return(&strgDto.GetUserByUsernameResponse{
		ID:           id,
		HashPassword: []byte{'1'},
	}, nil)
	mockUserStrgIntf.On("Insert", ctx, mock.Anything).Return(nil)

	response, err := authService.SignIn(ctx, req)

	assert.Equal(t, service.ErrIncorrectPassword, err)
	assert.Empty(t, response)
}

// empty username
func TestSignInFailed_02(t *testing.T) {

	mockUserStrgIntf := new(mocks.UserIntf)
	mockLogger := new(loggerMock.Interface)

	jwtMngIntf := jwt.NewJwtManager("avito", 1)

	authService := service.NewAuthSvc(mockLogger, jwtMngIntf, mockUserStrgIntf)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &svcDto.SignInRequest{
		Username: "",
		Password: "1",
	}

	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Warnf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Times(0)

	response, err := authService.SignIn(ctx, req)

	assert.Equal(t, service.ErrEmptyUsername, err)
	assert.Empty(t, response)
}

//empty password
func TestSignInFailed_03(t *testing.T) {

	mockUserStrgIntf := new(mocks.UserIntf)
	mockLogger := new(loggerMock.Interface)
	jwtMngIntf := jwt.NewJwtManager("avito", 1)

	authService := service.NewAuthSvc(mockLogger, jwtMngIntf, mockUserStrgIntf)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &svcDto.SignInRequest{
		Username: "1",
		Password: "",
	}

	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Warnf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Times(0)

	response, err := authService.SignIn(ctx, req)

	assert.Equal(t, service.ErrEmptyPassword, err)
	assert.Empty(t, response)
}
