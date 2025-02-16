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
	ErrSenderDoesntExist    = errors.New("receiver with this username doesnt exist")
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

    query := `
        select u1.id, u1.coins_amount, u2.id
        from users u1
        join users u2 ON u2.username = $2
        where u1.id = $1
        for update nowait`

    var (
        senderID uuid.UUID
        senderBalance int32
        receiverID uuid.UUID
    )

    err = tx.QueryRow(
        ctx,
        query,
        request.FromUserID,
        request.ToUsername,
    ).Scan(&senderID, &senderBalance, &receiverID)

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return ErrReceiverDoesntExist
        }
        return ErrSelectQueryRow
    }

    if senderID == receiverID {
        return ErrSameUser
    }

    if senderBalance < request.CoinsAmount {
        return ErrNotEnoughCoins
    }

    query = `
        update users 
        set coins_amount = CASE 
            when id = $1 THEN coins_amount - $3
            when id = $2 THEN coins_amount + $3
        end
        where id IN ($1, $2)`

    _, err = tx.Exec(ctx, query, senderID, receiverID, request.CoinsAmount)
    if err != nil {
        return ErrIncreaseCoinsAmount
    }

    query = `
        insert into transactions (from_user_id, to_user_id, coins_amount) 
        values ($1, $2, $3)`

    _, err = tx.Exec(ctx, query, senderID, receiverID, request.CoinsAmount)
    if err != nil {
        return ErrSenderDoesntExist
    }

    if err = tx.Commit(ctx); err != nil {
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
