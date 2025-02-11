package service

import (
	svcDto "avito-task-2025/backend/internal/service/dto"
	"avito-task-2025/backend/internal/storage"
	strgDto "avito-task-2025/backend/internal/storage/dto"
	"avito-task-2025/backend/pkg/logger"
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	errCoinsInvalidUserID  = errors.New("trying to send coins with invalid UserID")
	errToUsername          = errors.New("username, which should recieve coins, is invalid")
	errNegativeCoinsAmount = errors.New("coins amount less then 0")
)

type CoinIntf interface {
	Send(ctx context.Context, request *svcDto.SendCoinsRequest) (err error)
}

type CoinSvc struct {
	logger          logger.Interface
	transactionIntf storage.TransactionIntf
}

func NewCoinSvc(logger logger.Interface, transactionIntf storage.TransactionIntf) CoinIntf {
	return &CoinSvc{
		logger:          logger,
		transactionIntf: transactionIntf,
	}
}

func (s *CoinSvc) Send(ctx context.Context, request *svcDto.SendCoinsRequest) (err error) {

	_, err = uuid.Parse(request.UserID.String())
	if err != nil {
		s.logger.Errorf(errCoinsInvalidUserID.Error())
		return errCoinsInvalidUserID
	}

	if request.ToUserUsername == "" {
		s.logger.Warnf(errNegativeCoinsAmount.Error())
		return errNegativeCoinsAmount
	}

	if request.CoinsAmount < 0 {
		s.logger.Warnf(errToUsername.Error())
		return errToUsername
	}

	err = s.transactionIntf.Insert(ctx, &strgDto.InsertTransactionRequest{
		FromUserID:  request.UserID,
		ToUsername:  request.ToUserUsername,
		CoinsAmount: request.CoinsAmount,
	})
	if err != nil {
		s.logger.Infof(err.Error())
		return err
	}

	return nil
}
