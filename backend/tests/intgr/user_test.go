package intgr

import (
	"context"
	"testing"
	"time"

	"avito-task-2025/backend/tests"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	strgDto "avito-task-2025/backend/internal/storage/dto"
	"avito-task-2025/backend/internal/storage/postgres"
	// "avito-task-2025/backend/internal/storage/"
)

// success sign up
func TestInsertUserSuccess(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := tests.NewTestConfig("file://../../../db/postgres/test_migrations/intgr/user")
	_ = migrator.Force(01)
	_ = migrator.Down()

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	userStrgIntf := postgres.NewUserStrg(dbConnector)

	req := &strgDto.InsertUserRequest{
		Username:     "test1",
		HashPassword: []byte{'0'},
	}

	res, err := userStrgIntf.Insert(ctx, req)

	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}

func TestGetUserByUsernameSuccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := tests.NewTestConfig("file://../../../db/postgres/test_migrations/intgr/user")
	_ = migrator.Force(01)
	_ = migrator.Down()
	_ = migrator.Up()

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	userStrgIntf := postgres.NewUserStrg(dbConnector)

	req := &strgDto.GetUserByUsernameRequest{
		Username: "test",
	}

	res, err := userStrgIntf.GetByUsername(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, []byte{'9'}, res.HashPassword)
}

func TestGetUserByUsernameEmptyPassword(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := tests.NewTestConfig("file://../../../db/postgres/test_migrations/intgr/user")
	_ = migrator.Force(2)
	_ = migrator.Down()
	_ =  migrator.Up()

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	userStrgIntf := postgres.NewUserStrg(dbConnector)

	req := &strgDto.GetUserByUsernameRequest{
		Username: "testest",
	}

	res, err := userStrgIntf.GetByUsername(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, []byte(nil), res.HashPassword)
}

func TestGetCoinsByUserID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := tests.NewTestConfig("file://../../../db/postgres/test_migrations/intgr/user")
	_ = migrator.Force(3)
	_ = migrator.Down()
	_ = migrator.Up()

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	userStrgIntf := postgres.NewUserStrg(dbConnector)

	id, _ := uuid.Parse("068d0e53-9826-4fe0-8d86-b925c52ae25c")

	req := &strgDto.GetCoinsByUserIDRequest{
		UserID: id,
	}

	res, err := userStrgIntf.GetCoinsByUserID(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, int32(1000), res.Amount)
}

func TestGetCoinsByUserIDInvalidID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := tests.NewTestConfig("file://../../../db/postgres/test_migrations/intgr/user")
	_ = migrator.Force(4)
	_ = migrator.Down()
	_ = migrator.Up()

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	userStrgIntf := postgres.NewUserStrg(dbConnector)

	id, _ := uuid.Parse("068d0e53-9826-4fe0-8d86-b925c52ae25a")

	req := &strgDto.GetCoinsByUserIDRequest{
		UserID: id,
	}

	_, err = userStrgIntf.GetCoinsByUserID(ctx, req)

	assert.Equal(t, postgres.ErrInvalidUserID, err)
}
