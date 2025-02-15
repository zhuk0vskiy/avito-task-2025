package intgr

import (
	strgDto "avito-task-2025/backend/internal/storage/dto"
	"avito-task-2025/backend/internal/storage/postgres"
	"avito-task-2025/backend/tests"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInsertTransactionSuccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := tests.NewTestConfig("file://../../../db/postgres/test_migrations/intgr/transaction")
	_ = migrator.Force(01)
	_ = migrator.Down()
	_ = migrator.Up()

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

	assert.NoError(t, err)
}

func TestReceiverDoesntExist(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := tests.NewTestConfig("file://../../../db/postgres/test_migrations/intgr/transaction")
	_ = migrator.Force(01)
	_ = migrator.Down()
	_ = migrator.Up()

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

	assert.Equal(t, postgres.ErrReceiverDoesntExist, err)
}

func TestNotEnoughCoins(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := tests.NewTestConfig("file://../../../db/postgres/test_migrations/intgr/transaction")
	_ = migrator.Force(01)
	_ = migrator.Down()
	_ = migrator.Up()

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

	assert.Equal(t, postgres.ErrNotEnoughCoins, err)
}

func TestSenderDoesntExist(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := tests.NewTestConfig("file://../../../db/postgres/test_migrations/intgr/transaction")
	_ = migrator.Force(01)
	_ = migrator.Down()
	_ = migrator.Up()

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

	assert.Equal(t, postgres.ErrSenderDoesntExist, err)
}

func TestGetTransactionToUserID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := tests.NewTestConfig("file://../../../db/postgres/test_migrations/intgr/transaction")
	_ = migrator.Force(02)
	_ = migrator.Down()
	_ = migrator.Up()

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

	_ = migrator.Down()
	assert.NoError(t, err)
	assert.Equal(t, int32(101), res.Transactions[0].CoinsAmount)
	assert.Equal(t, "test2", res.Transactions[0].FromUsername)
}

func TestGetTransactionFromUserID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := tests.NewTestConfig("file://../../../db/postgres/test_migrations/intgr/transaction")
	_ = migrator.Force(02)
	_ = migrator.Down()
	_ = migrator.Up()

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	transStrgIntf := postgres.NewTransactionStrg(dbConnector)

	id, _ := uuid.Parse("068d0e53-9826-4fe0-8d86-b925c52ae25c")

	req := &strgDto.GetTransactionFromUserIDRequest{
		UserID: id,
	}

	res, err := transStrgIntf.GetFromUserID(ctx, req)

	_ = migrator.Down()
	assert.NoError(t, err)
	assert.Equal(t, int32(100), res.Transactions[0].CoinsAmount)
	assert.Equal(t, "test2", res.Transactions[0].ToUsername)
}
