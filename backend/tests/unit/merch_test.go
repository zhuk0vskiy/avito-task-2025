package unit

import (
	"avito-task-2025/backend/internal/service"
	"avito-task-2025/backend/internal/storage/mocks"
	loggerMock "avito-task-2025/backend/pkg/logger/mocks"
	"context"
	"testing"
	"time"

	svcDto "avito-task-2025/backend/internal/service/dto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Valid UUID and merch name should successfully insert bought merch record
func TestBuyMerchSuccess(t *testing.T) {
	mockBoughtMerchIntf := new(mocks.BoughtMerchIntf)
	mockLogger := new(loggerMock.Interface)

	merchSvcIntf := service.NewMerchSvc(mockLogger, mockBoughtMerchIntf)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, _ := uuid.NewRandom()
	merchName := "test_merch"
	req := &svcDto.BuyMerchRequest{
		UserID:    id,
		MerchName: merchName,
	}

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Warnf", mock.Anything, mock.Anything).Times(0)

	mockBoughtMerchIntf.On("Insert", mock.Anything, mock.Anything).Return(nil)

	err := merchSvcIntf.Buy(ctx, req)

	assert.NoError(t, err)
}

func TestBuyMerchInvalidUUID(t *testing.T) {
	mockBoughtMerchIntf := new(mocks.BoughtMerchIntf)
	mockLogger := new(loggerMock.Interface)

	merchSvcIntf := service.NewMerchSvc(mockLogger, mockBoughtMerchIntf)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	invalidUUID := uuid.Nil
	merchName := "test_merch"
	req := &svcDto.BuyMerchRequest{
		UserID:    invalidUUID,
		MerchName: merchName,
	}

	mockLogger.On("Errorf", mock.Anything).Return()
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Warnf", mock.Anything, mock.Anything).Times(0)

	mockBoughtMerchIntf.On("Insert", ctx, mock.Anything).Times(0)

	err := merchSvcIntf.Buy(ctx, req)

	assert.Equal(t, service.ErrMerchInvalidUserID, err)
}


func TestBuyMerchWithEmptyMerchName(t *testing.T) {

	mockBoughtMerchIntf := new(mocks.BoughtMerchIntf)
	mockLogger := new(loggerMock.Interface)

	merchSvc := service.NewMerchSvc(mockLogger, mockBoughtMerchIntf)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, _ := uuid.NewRandom()
	req := &svcDto.BuyMerchRequest{
		UserID:    id,
		MerchName: "",
	}

	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
	mockLogger.On("Warnf", service.ErrEmptyMerchName.Error()).Times(1)
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Times(0)

	err := merchSvc.Buy(ctx, req)

	assert.Equal(t, err, service.ErrEmptyMerchName)
}
