package unit

import (
	"avito-task-2025/backend/internal/entity"
	"avito-task-2025/backend/internal/service"
	svcDto "avito-task-2025/backend/internal/service/dto"
	strgDto "avito-task-2025/backend/internal/storage/dto"
	"avito-task-2025/backend/internal/storage/mocks"
	loggerMock "avito-task-2025/backend/pkg/logger/mocks"

	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Valid UUID returns complete user info with coins, inventory and transaction history
func TestGetInfoSuccess(t *testing.T) {
	mockUserIntf := new(mocks.UserIntf)
	mockBoughtMerchIntf := new(mocks.BoughtMerchIntf)
	mockTransactionIntf := new(mocks.TransactionIntf)
	mockLogger := new(loggerMock.Interface)

	userSvc := service.NewUserSvc(mockLogger, mockUserIntf, mockBoughtMerchIntf, mockTransactionIntf)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, _ := uuid.NewRandom()
	req := &svcDto.GetUserInfoRequest{
		UserID: id,
	}

	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)

	mockUserIntf.On("GetCoinsByUserID", ctx, &strgDto.GetCoinsByUserIDRequest{UserID: id}).
		Return(&strgDto.GetCoinsByUserIDResponse{}, nil)

	mockBoughtMerchIntf.On("GetByUserID", ctx, &strgDto.GetBoughtMerchByUserIDRequest{UserID: id}).
		Return(&strgDto.GetBoughtMerchByUserIDResponse{Merchs: []*entity.Merch{}}, nil)

	mockTransactionIntf.On("GetToUserID", ctx, &strgDto.GetTransactionToUserIDRequest{UserID: id}).
		Return(&strgDto.GetTransactionToUserIDResponse{Transactions: []*entity.Transaction{}}, nil)

	mockTransactionIntf.On("GetFromUserID", ctx, &strgDto.GetTransactionFromUserIDRequest{UserID: id}).
		Return(&strgDto.GetTransactionFromUserIDResponse{Transactions: []*entity.Transaction{}}, nil)

	response, err := userSvc.GetInfo(ctx, req)

	assert.NoError(t, err)

	assert.NotEmpty(t, response)
}

// Invalid UUID format in request returns ErrInfoInvalidUserID
func TestGetInfoInvalidUUID(t *testing.T) {
	mockUserIntf := new(mocks.UserIntf)
	mockBoughtMerchIntf := new(mocks.BoughtMerchIntf)
	mockTransactionIntf := new(mocks.TransactionIntf)
	mockLogger := new(loggerMock.Interface)

	userSvc := service.NewUserSvc(mockLogger, mockUserIntf, mockBoughtMerchIntf, mockTransactionIntf)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &svcDto.GetUserInfoRequest{
		UserID: uuid.Nil,
	}

	mockLogger.On("Errorf", mock.Anything).Return()

	response, err := userSvc.GetInfo(ctx, req)

	assert.Equal(t, service.ErrInfoInvalidUserID, err)
	assert.Nil(t, response)
}
