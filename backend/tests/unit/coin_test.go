package unit

import (
	"avito-task-2025/backend/internal/service"
	svcDto "avito-task-2025/backend/internal/service/dto"
	"avito-task-2025/backend/internal/storage/mocks"
	loggerMock "avito-task-2025/backend/pkg/logger/mocks"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSendCoinsSuccess(t *testing.T) {

	mockTransactionStrgIntf := new(mocks.TransactionIntf)
	mockLogger := new(loggerMock.Interface)

	coinSvcIntf := service.NewCoinSvc(mockLogger, mockTransactionStrgIntf)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, _ := uuid.NewRandom()
	toUsername := "test"
	coinsAmount := 1
	req := &svcDto.SendCoinsRequest{
		UserID:      id,
		ToUsername:  toUsername,
		CoinsAmount: int32(coinsAmount),
	}

	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Warnf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Times(0)

	mockTransactionStrgIntf.On("Insert", mock.Anything, mock.Anything).Return(nil)

	err := coinSvcIntf.Send(ctx, req)

	assert.NoError(t, err)
}

func TestSendCoinsInvalidUUID(t *testing.T) {
	mockTransactionStrgIntf := new(mocks.TransactionIntf)
	mockLogger := new(loggerMock.Interface)

	coinSvcIntf := service.NewCoinSvc(mockLogger, mockTransactionStrgIntf)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	invalidUUID := uuid.Nil
	req := &svcDto.SendCoinsRequest{
		UserID:      invalidUUID,
		ToUsername:  "testuser",
		CoinsAmount: 100,
	}

	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Warnf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Times(0)

	mockTransactionStrgIntf.On("Insert", mock.Anything, mock.Anything).Return(nil)

	err := coinSvcIntf.Send(ctx, req)

	assert.Equal(t, service.ErrCoinsInvalidUserID, err)
}

func TestSendCoinsEmptyUsername(t *testing.T) {

	mockTransactionStrgIntf := new(mocks.TransactionIntf)
	mockLogger := new(loggerMock.Interface)

	coinSvcIntf := service.NewCoinSvc(mockLogger, mockTransactionStrgIntf)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, _ := uuid.NewRandom()
	toUsername := ""
	coinsAmount := 1
	req := &svcDto.SendCoinsRequest{
		UserID:      id,
		ToUsername:  toUsername,
		CoinsAmount: int32(coinsAmount),
	}

	mockLogger.On("Warnf", service.ErrNegativeCoinsAmount.Error()).Once()

	err := coinSvcIntf.Send(ctx, req)

	assert.Equal(t, service.ErrNegativeCoinsAmount, err)
}

func TestSendCoinsZeroAmountReturnsErrInvalidToUsername(t *testing.T) {

	mockTransactionStrgIntf := new(mocks.TransactionIntf)
	mockLogger := new(loggerMock.Interface)

	coinSvcIntf := service.NewCoinSvc(mockLogger, mockTransactionStrgIntf)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, _ := uuid.NewRandom()
	toUsername := "test"
	coinsAmount := 0
	req := &svcDto.SendCoinsRequest{
		UserID:      id,
		ToUsername:  toUsername,
		CoinsAmount: int32(coinsAmount),
	}

	mockLogger.On("Warnf", service.ErrInvalidToUsername.Error()).Once()

	err := coinSvcIntf.Send(ctx, req)

	assert.Equal(t, service.ErrInvalidToUsername, err)
}
