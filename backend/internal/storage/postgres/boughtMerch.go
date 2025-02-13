package postgres

import (
	"avito-task-2025/backend/internal/entity"
	"avito-task-2025/backend/internal/storage"
	strgDto "avito-task-2025/backend/internal/storage/dto"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	errBuyTransStart           = errors.New("failed to start buy merch trans")
	errBuyTransRollback        = errors.New("failed to rollback buy merch trans")
	errBuyDecreaseCoinsAmount  = errors.New("user dont have enough coins")
	errBuyNotEnoughMoney       = errors.New("user dont have enough money to by merch")
	errBuyTransRecord          = errors.New("failed to insert bought merch record")
	errBuyCommitTrans          = errors.New("failed to commit buy merch trans")
	errGetAllBoughtMerchByUser = errors.New("falied to get user bought merchs")
	errScanMerchRow            = errors.New("failed to scan bought merch row")
	errInvalidMerchName        = errors.New("this merch doesnt exist")
	errIncreaseBoughtMerchAmount     = errors.New("failed to increase bought merch amount")
)

var (
	Test_errBuyDecreaseCoinsAmount = errBuyDecreaseCoinsAmount
	Test_errBuyTransRecord         = errBuyTransRecord
)

type BoughtMerchStrg struct {
	dbConnector *pgxpool.Pool
}

func NewBoughtMerchStrg(dbConnector *pgxpool.Pool) storage.BoughtMerchIntf {
	return &BoughtMerchStrg{
		dbConnector: dbConnector,
	}
}

func (s *BoughtMerchStrg) Insert(ctx context.Context, request *strgDto.AddBoughtMerchRequest) (err error) {


	var cost int
	var merchID uuid.UUID

	query := `select id, cost from merchs where type = $1`
	err = s.dbConnector.QueryRow(
		ctx,
		query,
		request.MerchName,
	).Scan(
		&merchID,
		&cost,
	)
	if err != nil {
		return errInvalidMerchName
	}

	tx, err := s.dbConnector.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return errBuyTransStart
	}
	defer func() {
		if err != nil {
			rollbackerr := tx.Rollback(ctx)
			if rollbackerr != nil {
				err = errBuyTransRollback
			}
		}
	}()

	var coinsAmount int32
	query = `
			update users 
			set coins_amount = coins_amount - $1
			where id = $2
			returning coins_amount`
	err = tx.QueryRow(
		ctx,
		query,
		cost,
		request.UserID,
	).Scan(
		&coinsAmount,
	)
	if err != nil {
		return errBuyDecreaseCoinsAmount
	}
	if coinsAmount < 0 {
		return errBuyNotEnoughMoney
	}


	var tmp *int
	var merchsAmount int
	fmt.Println(request.UserID, merchID)
	query = `insert into bought_merchs(user_id, merch_id) 
			values ($1, $2)
			returning amount`
	err = tx.QueryRow(
		ctx,
		query,
		request.UserID,
		merchID,
	).Scan(
		&tmp,
	)

	if err != nil {
		return errBuyTransRecord
	}

	if tmp != nil {
		merchsAmount = *tmp
	} else {
		merchsAmount = 0
	}

	merchsAmount++
	query = `update bought_merchs set amount = $1`
	_, err = tx.Exec(
		ctx,
		query,
		merchsAmount,
	)
	if err != nil {
		return errIncreaseBoughtMerchAmount
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errBuyCommitTrans
	}

	return nil
}

func (s *BoughtMerchStrg) GetByUserID(ctx context.Context, request *strgDto.GetBoughtMerchByUserIDRequest) (response *strgDto.GetBoughtMerchByUserIDResponse, err error) {
	query := `select bought_merchs.amount, merchs.type from bought_merchs join merchs on bought_merchs.merch_id = merchs.id where bought_merchs.user_id = $1`

	rows, err := s.dbConnector.Query(
		ctx,
		query,
		request.UserID,
	)
	if err != nil {
		return nil, errGetAllBoughtMerchByUser
	}
	defer rows.Close()

	merchs := make([]*entity.Merch, 0)

	for rows.Next() {
		merch := entity.Merch{}
		err = rows.Scan(
			&merch.Amount,
			&merch.Type,
		)
		if err != nil {
			return nil, errScanMerchRow
		}

		merchs = append(merchs, &merch)
	}

	response = &strgDto.GetBoughtMerchByUserIDResponse{
		Merchs: merchs,
	}
	return response, nil
}
