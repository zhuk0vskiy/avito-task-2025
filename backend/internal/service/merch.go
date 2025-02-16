package service

import (
	svcDto "avito-task-2025/backend/internal/service/dto"
	"avito-task-2025/backend/internal/storage"
	strgDto "avito-task-2025/backend/internal/storage/dto"
	"avito-task-2025/backend/pkg/logger"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEmptyMerchName     = errors.New("merch name is empty")
	ErrMerchInvalidUserID = errors.New("trying to buy merch with invalid user id")
)

type MerchIntf interface {
	Buy(ctx context.Context, request *svcDto.BuyMerchRequest) (err error)
}

type MerchSvc struct {
	logger          logger.Interface
	boughtMerchIntf storage.BoughtMerchIntf
}

func NewMerchSvc(logger logger.Interface, boughtMerchIntf storage.BoughtMerchIntf) MerchIntf {
	return &MerchSvc{
		logger:          logger,
		boughtMerchIntf: boughtMerchIntf,
	}
}

func (s *MerchSvc) Buy(ctx context.Context, request *svcDto.BuyMerchRequest) (err error) {
	_, err = uuid.Parse(request.UserID.String())
	if err != nil || request.UserID == uuid.Nil {
		s.logger.Warnf(ErrMerchInvalidUserID.Error())
		return ErrMerchInvalidUserID
	}

	if request.MerchName == "" {
		s.logger.Warnf(ErrEmptyMerchName.Error())
		return ErrEmptyMerchName
	}

	s.logger.Debugf("before buy insert %s", time.Now().UnixMilli())
	err = s.boughtMerchIntf.Insert(ctx, &strgDto.InsertBoughtMerchRequest{
		UserID: request.UserID,
		Type:   request.MerchName,
	})
	s.logger.Debugf("after buy insert %s", time.Now().UnixMilli())
	if err != nil {
		s.logger.Infof(err.Error())
		return err
	}

	return nil
}
