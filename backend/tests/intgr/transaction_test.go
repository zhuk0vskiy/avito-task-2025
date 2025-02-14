package intgr

import (
	strgDto "avito-task-2025/backend/internal/storage/dto"
	"avito-task-2025/backend/internal/storage/postgres"
	"avito-task-2025/backend/tests/helper"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// success send
func TestInsertTransactionSuccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := helper.NewTestConfig("file://../../../db/postgres/test_migrations/transaction")
	migrator.Force(01)
	migrator.Down()
	migrator.Up()

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	transStrgIntf := postgres.NewTransactionStrg(dbConnector)

	id, _ := uuid.Parse("068d0e53-9826-4fe0-8d86-b925c52ae25c")

	req := &strgDto.InsertTransactionRequest{
		FromUserID:  id,
		ToUsername:  "test2",
		CoinsAmount: 10,
	}

	err = transStrgIntf.Insert(ctx, req)

	migrator.Down()
	assert.NoError(t, err)
}

func TestReceiverDoesntExist(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := helper.NewTestConfig("file://../../../db/postgres/test_migrations/transaction")
	migrator.Force(01)
	migrator.Down()
	migrator.Up()

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	transStrgIntf := postgres.NewTransactionStrg(dbConnector)

	id, _ := uuid.Parse("068d0e53-9826-4fe0-8d86-b925c52ae25c")

	req := &strgDto.InsertTransactionRequest{
		FromUserID:  id,
		ToUsername:  "test3",
		CoinsAmount: 10,
	}

	err = transStrgIntf.Insert(ctx, req)

	migrator.Down()
	assert.Equal(t, postgres.ErrReceiverDoesntExist, err)
}

func TestNotEnoughCoins(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := helper.NewTestConfig("file://../../../db/postgres/test_migrations/transaction")
	migrator.Force(01)
	migrator.Down()
	migrator.Up()

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	transStrgIntf := postgres.NewTransactionStrg(dbConnector)

	id, _ := uuid.Parse("068d0e53-9826-4fe0-8d86-b925c52ae25c")

	req := &strgDto.InsertTransactionRequest{
		FromUserID:  id,
		ToUsername:  "test2",
		CoinsAmount: 1001,
	}

	err = transStrgIntf.Insert(ctx, req)

	migrator.Down()
	assert.Equal(t, postgres.ErrNotEnoughCoins, err)
}

func TestSenderDoesntExist(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := helper.NewTestConfig("file://../../../db/postgres/test_migrations/transaction")
	migrator.Force(01)
	migrator.Down()
	migrator.Up()

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	transStrgIntf := postgres.NewTransactionStrg(dbConnector)

	id, _ := uuid.Parse("068d0e53-9826-4fe0-8d86-b925c52ae25b")

	req := &strgDto.InsertTransactionRequest{
		FromUserID:  id,
		ToUsername:  "test2",
		CoinsAmount: 1,
	}

	err = transStrgIntf.Insert(ctx, req)

	migrator.Down()
	assert.Equal(t, postgres.ErrSenderDoesntExist, err)
}

func TestGetTransactionToUserID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := helper.NewTestConfig("file://../../../db/postgres/test_migrations/transaction")
	migrator.Force(02)
	migrator.Down()
	migrator.Up()

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	transStrgIntf := postgres.NewTransactionStrg(dbConnector)

	id, _ := uuid.Parse("068d0e53-9826-4fe0-8d86-b925c52ae25c")

	req := &strgDto.GetTransactionToUserIDRequest{
		UserID: id,
	}

	res, err := transStrgIntf.GetToUserID(ctx, req)

	migrator.Down()
	assert.NoError(t, err)
	assert.Equal(t, int32(100), res.Transactions[0].CoinsAmount)
	assert.Equal(t, "test2", res.Transactions[0].FromUsername)
}

func TestGetTransactionFromUserID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := helper.NewTestConfig("file://../../../db/postgres/test_migrations/transaction")
	migrator.Force(02)
	migrator.Down()
	migrator.Up()

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	transStrgIntf := postgres.NewTransactionStrg(dbConnector)

	id, _ := uuid.Parse("068d0e53-9826-4fe0-8d86-b925c52ae25a")

	req := &strgDto.GetTransactionToUserIDRequest{
		UserID: id,
	}

	res, err := transStrgIntf.GetToUserID(ctx, req)

	migrator.Down()
	assert.NoError(t, err)
	assert.Equal(t, int32(100), res.Transactions[0].CoinsAmount)
	assert.Equal(t, "test1", res.Transactions[0].FromUsername)
}
