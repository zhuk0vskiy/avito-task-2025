package intgr

import (
	"context"
	"testing"
	"time"

	"avito-task-2025/backend/tests/helper"

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

	migrator, cfg := helper.NewTestConfig("file://../../../db/postgres/test_migrations/user")
	migrator.Force(00)

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	userStrgIntf := postgres.NewUserStrg(dbConnector)

	req := &strgDto.InsertUserRequest{
		Username:     "test1",
		HashPassword: []byte{'0'},
	}

	err = userStrgIntf.Insert(ctx, req)

	migrator.Down()
	assert.NoError(t, err)
}

func TestGetUserByUsernameSuccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := helper.NewTestConfig("file://../../../db/postgres/test_migrations/user")
	migrator.Force(001)
	migrator.Down()
	migrator.Up()

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	userStrgIntf := postgres.NewUserStrg(dbConnector)

	req := &strgDto.GetUserByUsernameRequest{
		Username: "test",
	}

	res, err := userStrgIntf.GetByUsername(ctx, req)

	migrator.Down()

	assert.NoError(t, err)
	assert.Equal(t, []byte{'9'}, res.HashPassword)
}

func TestGetUserByUsernameEmptyPassword(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := helper.NewTestConfig("file://../../../db/postgres/test_migrations/user")
	migrator.Force(002)
	migrator.Down()
	migrator.Up()

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	userStrgIntf := postgres.NewUserStrg(dbConnector)

	req := &strgDto.GetUserByUsernameRequest{
		Username: "testest",
	}

	res, err := userStrgIntf.GetByUsername(ctx, req)

	migrator.Down()

	assert.NoError(t, err)
	assert.Equal(t, []byte(nil), res.HashPassword)
}

func TestGetCoinsByUserID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := helper.NewTestConfig("file://../../../db/postgres/test_migrations/user")
	migrator.Down()
	migrator.Force(003)
	migrator.Up()

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

	migrator.Down()

	assert.NoError(t, err)
	assert.Equal(t, int32(1000), res.Amount)
}

func TestGetCoinsByUserIDInvalidID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := helper.NewTestConfig("file://../../../db/postgres/test_migrations/user")
	migrator.Force(004)
	migrator.Down()
	migrator.Up()

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

	migrator.Down()

	assert.Equal(t, postgres.ErrInvalidUserID, err)
}