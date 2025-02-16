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
	ErrCoinsInvalidUserID  = errors.New("trying to send coins with invalid UserID")
	ErrInvalidToUsername   = errors.New("username, which should recieve coins, is invalid")
	ErrNegativeCoinsAmount = errors.New("coins amount less then 1")
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
	if err != nil || request.UserID == uuid.Nil{
		s.logger.Warnf(ErrCoinsInvalidUserID.Error())
		return ErrCoinsInvalidUserID
	}

	if request.ToUsername == "" {
		s.logger.Warnf(ErrNegativeCoinsAmount.Error())
		return ErrNegativeCoinsAmount
	}

	if request.CoinsAmount <= 0 {
		s.logger.Warnf(ErrInvalidToUsername.Error())
		return ErrInvalidToUsername
	}

	err = s.transactionIntf.Insert(ctx, &strgDto.InsertTransactionRequest{
		FromUserID:  request.UserID,
		ToUsername:  request.ToUsername,
		CoinsAmount: request.CoinsAmount,
	})
	if err != nil {
		s.logger.Warnf(err.Error())
		return err
	}

	return nil
}
