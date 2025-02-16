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

func TestInsertBoughtMerchSuccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := tests.NewTestConfig("file://../../../db/postgres/test_migrations/intgr/bought_merch")

	_ = migrator.Force(01)
	_ = migrator.Down()
	_ = migrator.Migrate(01)


	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	boughtMerchStrgIntf := postgres.NewBoughtMerchStrg(dbConnector)

	id, _ := uuid.Parse("068d0e53-9826-4fe0-8d86-b925c52ae25c")

	req := &strgDto.InsertBoughtMerchRequest{
		UserID: id,
		Type:   "t-shirt",
	}

	err = boughtMerchStrgIntf.Insert(ctx, req)

	assert.NoError(t, err)
}

func Test2InsertsBoughtMerchSuccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := tests.NewTestConfig("file://../../../db/postgres/test_migrations/intgr/bought_merch")

	_ = migrator.Force(01)
	_ = migrator.Down()
	_ = migrator.Migrate(01)

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	boughtMerchStrgIntf := postgres.NewBoughtMerchStrg(dbConnector)

	id, _ := uuid.Parse("068d0e53-9826-4fe0-8d86-b925c52ae25c")

	req := &strgDto.InsertBoughtMerchRequest{
		UserID: id,
		Type:   "t-shirt",
	}

	_ = boughtMerchStrgIntf.Insert(ctx, req)
	err = boughtMerchStrgIntf.Insert(ctx, req)

	assert.NoError(t, err)
}

func TestInvalidMerchName(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := tests.NewTestConfig("file://../../../db/postgres/test_migrations/intgr/bought_merch")

	_ = migrator.Force(01)
	_ = migrator.Down()
	_ = migrator.Migrate(01)

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	boughtMerchStrgIntf := postgres.NewBoughtMerchStrg(dbConnector)

	id, _ := uuid.Parse("068d0e53-9826-4fe0-8d86-b925c52ae25c")

	req := &strgDto.InsertBoughtMerchRequest{
		UserID: id,
		Type:   "avito",
	}

	err = boughtMerchStrgIntf.Insert(ctx, req)

	assert.Equal(t, postgres.ErrInvalidMerchName, err)
}

func TestNotEnoughCoinsToBuy(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := tests.NewTestConfig("file://../../../db/postgres/test_migrations/intgr/bought_merch")

	_ = migrator.Migrate(02)

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	boughtMerchStrgIntf := postgres.NewBoughtMerchStrg(dbConnector)

	id, _ := uuid.Parse("068d0e53-9826-4fe0-8d86-b925c52ae25c")

	req := &strgDto.InsertBoughtMerchRequest{
		UserID: id,
		Type:   "t-shirt",
	}

	err = boughtMerchStrgIntf.Insert(ctx, req)

	assert.Equal(t, postgres.ErrBuyNotEnoughCoins, err)
}

func TestBoughtMerchByUserID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := tests.NewTestConfig("file://../../../db/postgres/test_migrations/intgr/bought_merch")

	_ = migrator.Force(3)
	_ = migrator.Down()
	_ = migrator.Up()

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	boughtMerchStrgIntf := postgres.NewBoughtMerchStrg(dbConnector)

	id, _ := uuid.Parse("068d0e53-9826-4fe0-8d86-b925c52ae25c")

	req := &strgDto.GetBoughtMerchByUserIDRequest{
		UserID: id,
	}

	res, err := boughtMerchStrgIntf.GetByUserID(ctx, req)

	assert.NotEmpty(t, res)
	assert.NoError(t, err)

}
