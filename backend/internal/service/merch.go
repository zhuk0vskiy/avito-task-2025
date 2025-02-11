package service

import (
	svcDto "avito-task-2025/backend/internal/service/dto"
	strgDto "avito-task-2025/backend/internal/storage/dto"
	"avito-task-2025/backend/internal/storage"
	"avito-task-2025/backend/pkg/logger"
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	errEmptyMerchName = errors.New("merch name is emtpy")
	errMerchInvalidUserID = errors.New("trying to buy merch with invalid UserID")
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
	if err != nil {
		s.logger.Errorf(errMerchInvalidUserID.Error())
		return errMerchInvalidUserID
	}

	if request.MerchName == "" {
		s.logger.Warnf(errEmptyMerchName.Error())
		return errEmptyMerchName
	}

	err = s.boughtMerchIntf.Insert(ctx, &strgDto.AddBoughtMerchRequest{
		UserID: request.UserID,
		MerchName: request.MerchName,
	})
	if err != nil {
		s.logger.Infof(err.Error())
		return err
	}

	return nil
}
