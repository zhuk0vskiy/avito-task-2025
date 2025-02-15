package postgres

import (
	"avito-task-2025/backend/internal/entity"
	"avito-task-2025/backend/internal/storage"
	strgDto "avito-task-2025/backend/internal/storage/dto"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInsertTransStart     = errors.New("failed to start send coins trans")
	ErrInsertTransRollback  = errors.New("failed to rollback send coins trans")
	ErrDecreaseCoinsAmount  = errors.New("failed to decrease user coins amount to send")
	ErrNotEnoughCoins       = errors.New("not enough coins to send")
	ErrIncreaseCoinsAmount  = errors.New("failed to increase user coins amount")
	ErrNoReceiveUser        = errors.New("receive user doesnt exist")
	ErrSenderDoesntExist    = errors.New("failed to insert send coins trans record")
	ErrInsertCommitTrans    = errors.New("failed to commit send coins trans")
	ErrGetTransNoSuchUserID = errors.New("no such user id")
	ErrSelectQueryRow       = errors.New("error while quering select row")
	ErrGetTransToUser       = errors.New("failed to get trans to user")

	ErrReceiverDoesntExist = errors.New("receiver with this username doesnt exist")
	ErrSameUser            = errors.New("cant send coins to yourself")
)

type TransactionStrg struct {
	dbConnector *pgxpool.Pool
}

func NewTransactionStrg(dbConnector *pgxpool.Pool) storage.TransactionIntf {
	return &TransactionStrg{
		dbConnector: dbConnector,
	}
}

func (s *TransactionStrg) Insert(ctx context.Context, request *strgDto.InsertTransactionRequest) (err error) {

	var receiverID uuid.UUID
	query := `select id from users where username = $1`
	err = s.dbConnector.QueryRow(
		ctx,
		query,
		request.ToUsername,
	).Scan(
		&receiverID,
	)
	if err != nil {
		return ErrReceiverDoesntExist
	}
	if receiverID == request.FromUserID {
		return ErrSameUser
	}

	tx, err := s.dbConnector.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return ErrInsertTransStart
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = ErrInsertTransRollback
			}
		}
	}()

	var coinsAmount int32
	query = `update users set coins_amount = coins_amount - $1 where id = $2 returning coins_amount`
	err = tx.QueryRow(
		ctx,
		query,
		request.CoinsAmount,
		request.FromUserID,
	).Scan(
		&coinsAmount,
	)

	if err != nil {
		return ErrSenderDoesntExist
	}
	if coinsAmount < 0 {
		return ErrNotEnoughCoins
	}

	var toUserID uuid.UUID
	query = `
        update users 
        set coins_amount = coins_amount + $1
        where id = $2
        returning coins_amount, id`
	err = tx.QueryRow(
		ctx,
		query,
		request.CoinsAmount,
		receiverID,
	).Scan(
		&coinsAmount,
		&toUserID,
	)
	if err != nil {
		return ErrIncreaseCoinsAmount
	}

	query = `insert into transactions(from_user_id, to_user_id, coins_amount) values ($1, $2, $3)`
	_, err = tx.Exec(
		ctx,
		query,
		request.FromUserID,
		toUserID,
		request.CoinsAmount,
	)
	if err != nil {
		return ErrSenderDoesntExist
	}

	err = tx.Commit(ctx)
	if err != nil {
		return ErrInsertCommitTrans
	}

	return nil
}

func (s *TransactionStrg) GetToUserID(ctx context.Context, request *strgDto.GetTransactionToUserIDRequest) (response *strgDto.GetTransactionToUserIDResponse, err error) {
	query := `select users.username, transactions.coins_amount from transactions join users on transactions.from_user_id = users.id where transactions.to_user_id = $1`

	rows, err := s.dbConnector.Query(
		ctx,
		query,
		request.UserID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrGetTransNoSuchUserID
		}
		return nil, ErrGetTransToUser
	}
	defer rows.Close()

	transactions := make([]*entity.Transaction, 0)
	for rows.Next() {
		transaction := entity.Transaction{}
		err = rows.Scan(
			&transaction.FromUsername,
			&transaction.CoinsAmount,
		)
		if err != nil {

			return nil, ErrSelectQueryRow
		}
		transactions = append(transactions, &transaction)
	}

	response = &strgDto.GetTransactionToUserIDResponse{
		Transactions: transactions,
	}

	return response, nil
}

func (s *TransactionStrg) GetFromUserID(ctx context.Context, request *strgDto.GetTransactionFromUserIDRequest) (response *strgDto.GetTransactionFromUserIDResponse, err error) {
	query := `select users.username, transactions.coins_amount from transactions join users on transactions.to_user_id = users.id where transactions.from_user_id = $1`

	rows, err := s.dbConnector.Query(
		ctx,
		query,
		request.UserID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrGetTransNoSuchUserID
		}
		return nil, ErrGetTransToUser
	}
	defer rows.Close()

	transactions := make([]*entity.Transaction, 0)
	for rows.Next() {
		transaction := entity.Transaction{}
		err = rows.Scan(
			&transaction.ToUsername,
			&transaction.CoinsAmount,
		)
		if err != nil {
			return nil, ErrSelectQueryRow
		}
		transactions = append(transactions, &transaction)
	}

	response = &strgDto.GetTransactionFromUserIDResponse{
		Transactions: transactions,
	}

	return response, nil
}
