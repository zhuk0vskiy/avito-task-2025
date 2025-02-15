package postgres

import (
	"context"
	"errors"

	"avito-task-2025/backend/internal/storage"
	strgDto "avito-task-2025/backend/internal/storage/dto"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrAddUser         = errors.New("failed to add user")
	ErrGetUserPassword = errors.New("failed to get user password")
	ErrGetCoinsAmount  = errors.New("failed to get user coins amount")
	ErrInvalidUserID   = errors.New("invlid user id to get coins")
)

type UserStrg struct {
	dbConnector *pgxpool.Pool
}

func NewUserStrg(dbConnector *pgxpool.Pool) storage.UserIntf {
	return &UserStrg{
		dbConnector: dbConnector,
	}
}

func (s *UserStrg) Insert(ctx context.Context, request *strgDto.InsertUserRequest) (response *strgDto.InsertUserResponse, err error) {
	response = &strgDto.InsertUserResponse{}
	query := `insert into users(username, password, coins_amount) values ($1, $2, 1000) returning id`

	err = s.dbConnector.QueryRow(
		ctx,
		query,
		request.Username,
		request.HashPassword,
	).Scan(
		&response.ID,
	)
	if err != nil {
		return nil, ErrAddUser
	}
	return response, nil
}

func (s *UserStrg) GetByUsername(ctx context.Context, request *strgDto.GetUserByUsernameRequest) (response *strgDto.GetUserByUsernameResponse, err error) {
	query := `select
				case
        			when exists (select 1 from users where username = $1)
        			then password
        			else '\x'::bytea
    			end as password, id 
             from users 
             where username = $1`

	response = &strgDto.GetUserByUsernameResponse{}

	err = s.dbConnector.QueryRow(
		ctx,
		query,
		request.Username,
	).Scan(
		&response.HashPassword,
		&response.ID,
	)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrGetUserPassword
		}
	}

	return response, nil
}

func (s *UserStrg) GetCoinsByUserID(ctx context.Context, request *strgDto.GetCoinsByUserIDRequest) (response *strgDto.GetCoinsByUserIDResponse, err error) {
	query := `select coins_amount from users where id = $1`

	response = &strgDto.GetCoinsByUserIDResponse{}
	err = s.dbConnector.QueryRow(
		ctx,
		query,
		request.UserID,
	).Scan(
		&response.Amount,
	)
	if err == pgx.ErrNoRows {
		return nil, ErrInvalidUserID
	}
	if err != nil {
		return nil, ErrGetCoinsAmount
	}
	// response.Amount = *tmp

	return response, nil
}
