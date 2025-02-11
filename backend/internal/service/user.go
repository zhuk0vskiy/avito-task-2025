package service

import (
	"avito-task-2025/backend/internal/entity"
	svcDto "avito-task-2025/backend/internal/service/dto"
	"avito-task-2025/backend/internal/storage"
	strgDto "avito-task-2025/backend/internal/storage/dto"
	"avito-task-2025/backend/pkg/logger"
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	errInfoInvalidUserID = errors.New("trying to get user info with invalid UserID")
)

type UserIntf interface {
	GetInfo(ctx context.Context, request *svcDto.GetUserInfoRequest) (response *svcDto.GetUserInfoResponse, err error)
}

type UserSvc struct {
	logger          logger.Interface
	userIntf        storage.UserIntf
	boughtMerchIntf storage.BoughtMerchIntf
	transactionIntf storage.TransactionIntf
}

func NewUserSvc(logger logger.Interface, userIntf storage.UserIntf, boughtMerchIntf storage.BoughtMerchIntf, transactionIntf storage.TransactionIntf) UserIntf {
	return &UserSvc{
		logger:          logger,
		userIntf:        userIntf,
		boughtMerchIntf: boughtMerchIntf,
		transactionIntf: transactionIntf,
	}
}

func (s *UserSvc) GetInfo(ctx context.Context, request *svcDto.GetUserInfoRequest) (response *svcDto.GetUserInfoResponse, err error) {
	_, err = uuid.Parse(request.UserID.String())
	if err != nil {
		s.logger.Errorf(errInfoInvalidUserID.Error())
		return nil, errInfoInvalidUserID
	}

	coins, err := s.userIntf.GetCoinsByUserID(ctx, &strgDto.GetCoinsByUserIDRequest{
		UserID: request.UserID,
	})
	if err != nil {
		s.logger.Infof(err.Error())
		return nil, err
	}

	inventory, err := s.boughtMerchIntf.GetByUserID(ctx, &strgDto.GetBoughtMerchByUserIDRequest{
		UserID: request.UserID,
	})
	if err != nil {
		s.logger.Infof(err.Error())
		return nil, err
	}

	receivedCoins, err := s.transactionIntf.GetByFromUserID(ctx, &strgDto.GetTransactionByFromUserIDRequest{
		UserID: request.UserID,
	})
	if err != nil {
		s.logger.Infof(err.Error())
		return nil, err
	}

	sentCoins, err := s.transactionIntf.GetByToUserID(ctx, &strgDto.GetTransactionByToUserIDRequest{
		UserID: request.UserID,
	})
	if err != nil {
		s.logger.Infof(err.Error())
		return nil, err
	}


	coinHistory := struct {
		Received []*entity.Transaction
		Sent     []*entity.Transaction
	}{
		Received: receivedCoins.Transactions,
		Sent: sentCoins.Transactions,
	}

	userInfo := struct {
		Coins       int32
		Inventory   []*entity.Merch
		CoinHistory *struct {
			Received []*entity.Transaction
			Sent     []*entity.Transaction
		}
	}{
		Coins: coins.Amount,
		Inventory: inventory.Merchs,
		CoinHistory: &coinHistory,
	}

	response = &svcDto.GetUserInfoResponse{
		UserInfo: &userInfo,
	}

	return response, nil

}
