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
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	errBuyTransStart           = errors.New("failed to start buy merch trans")
	errBuyTransRollback        = errors.New("failed to rollback buy merch trans")
	errBuyDecreaseCoinsAmount  = errors.New("failed to decrease coins")
	ErrBuyNotEnoughCoins       = errors.New("not enough coins to buy merch")
	errBuyTransRecord          = errors.New("failed to insert bought merch record")
	errBuyCommitTrans          = errors.New("failed to commit buy merch trans")
	errGetAllBoughtMerchByUser = errors.New("falied to get user bought merchs")
	errScanMerchRow            = errors.New("failed to scan bought merch row")
	ErrInvalidMerchName        = errors.New("this merch doesnt exist")
	// errIncreaseBoughtMerchAmount = errors.New("failed to increase bought merch amount")
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

func (s *BoughtMerchStrg) Insert(ctx context.Context, request *strgDto.InsertBoughtMerchRequest) (err error) {
	tx, err := s.dbConnector.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return errBuyTransStart
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				err = errBuyTransRollback
			}
		}
	}()

	var (
		merchID uuid.UUID
		cost    int32
	)

	query := `
        select id, cost 
        from merchs 
        where type = $1 
        for update`

	err = tx.QueryRow(ctx, query, request.Type).Scan(&merchID, &cost)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInvalidMerchName
		}
		return ErrSelectQueryRow
	}

	var userCoins int32
	query = `
        select coins_amount 
        from users 
        where id = $1 
        for update`

	err = tx.QueryRow(ctx, query, request.UserID).Scan(&userCoins)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInvalidUserID
		}
		return ErrSelectQueryRow
	}

	if userCoins < cost {
		return ErrBuyNotEnoughCoins
	}

	// Обновляем баланс пользователя
	query = `
        update users 
        set coins_amount = coins_amount - $1 
        where id = $2`

	_, err = tx.Exec(ctx, query, cost, request.UserID)
	if err != nil {
		return errBuyDecreaseCoinsAmount
	}

	query = `
    insert into bought_merchs (user_id, merch_id, amount) 
    values ($1, $2, 1) 
    on conflict (user_id, merch_id) 
    do update 
    set amount = EXCLUDED.amount + bought_merchs.amount
    returning amount
	`

	var newAmount int32
	err = tx.QueryRow(ctx, query, request.UserID, merchID).Scan(&newAmount)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return fmt.Errorf("postgres error: %s, detail: %s, where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
		}
		return fmt.Errorf("failed to update bought_merchs: %w", err)
	}
	if err = tx.Commit(ctx); err != nil {
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
			// fmt.Println(err)
			return nil, errScanMerchRow
		}

		merchs = append(merchs, &merch)
	}

	response = &strgDto.GetBoughtMerchByUserIDResponse{
		Merchs: merchs,
	}
	return response, nil
}
