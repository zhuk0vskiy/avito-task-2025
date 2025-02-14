package intgr

import (
	strgDto "avito-task-2025/backend/internal/storage/dto"
	"avito-task-2025/backend/internal/storage/postgres"
	"avito-task-2025/backend/tests/helper"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInsertBoughtMerchSuccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := helper.NewTestConfig("file://../../../db/postgres/test_migrations/bought_merch")

	migrator.Force(01)
	migrator.Down()
	migrator.Migrate(01)
	// migrator.Migrate(01)

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	boughtMerchsStrgIntf := postgres.NewBoughtMerchStrg(dbConnector)

	id, _ := uuid.Parse("068d0e53-9826-4fe0-8d86-b925c52ae25c")

	req := &strgDto.InsertBoughtMerchRequest{
		UserID:  id,
		Type:  "t-shirt",
	}

	err = boughtMerchsStrgIntf.Insert(ctx, req)

	migrator.Down()
	assert.NoError(t, err)
}

func Test2InsertsBoughtMerchSuccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := helper.NewTestConfig("file://../../../db/postgres/test_migrations/bought_merch")

	migrator.Force(01)
	migrator.Down()
	migrator.Migrate(01)

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	boughtMerchsStrgIntf := postgres.NewBoughtMerchStrg(dbConnector)

	id, _ := uuid.Parse("068d0e53-9826-4fe0-8d86-b925c52ae25c")

	req := &strgDto.InsertBoughtMerchRequest{
		UserID:  id,
		Type:  "t-shirt",
	}

	_ = boughtMerchsStrgIntf.Insert(ctx, req)
	err = boughtMerchsStrgIntf.Insert(ctx, req)

	// migrator.Down()
	assert.NoError(t, err)
}

func TestInvalidMerchName(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := helper.NewTestConfig("file://../../../db/postgres/test_migrations/bought_merch")

	migrator.Down()
	migrator.Migrate(01)

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	boughtMerchsStrgIntf := postgres.NewBoughtMerchStrg(dbConnector)

	id, _ := uuid.Parse("068d0e53-9826-4fe0-8d86-b925c52ae25c")

	req := &strgDto.InsertBoughtMerchRequest{
		UserID:  id,
		Type:  "avito",
	}

	err = boughtMerchsStrgIntf.Insert(ctx, req)

	assert.Equal(t, postgres.ErrInvalidMerchName, err)
}

func TestNotEnoughCoinsToBuy(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := helper.NewTestConfig("file://../../../db/postgres/test_migrations/bought_merch")
	// migrator.Migrate(01)
	migrator.Down()
	migrator.Migrate(02)

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	boughtMerchsStrgIntf := postgres.NewBoughtMerchStrg(dbConnector)

	id, _ := uuid.Parse("068d0e53-9826-4fe0-8d86-b925c52ae25c")

	req := &strgDto.InsertBoughtMerchRequest{
		UserID:  id,
		Type:  "t-shirt",
	}

	err = boughtMerchsStrgIntf.Insert(ctx, req)

	assert.Equal(t, postgres.ErrBuyNotEnoughCoins, err)
}

func TestBoughtMerchByUserID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, cfg := helper.NewTestConfig("file://../../../db/postgres/test_migrations/bought_merch")
	migrator.Migrate(01)
	// migrator.Migrate(03)

	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
	if err != nil {
		t.Error(err)
	}

	boughtMerchsStrgIntf := postgres.NewBoughtMerchStrg(dbConnector)

	id, _ := uuid.Parse("068d0e53-9826-4fe0-8d86-b925c52ae25c")

	req := &strgDto.GetBoughtMerchByUserIDRequest{
		UserID:  id,
	}

	res, err := boughtMerchsStrgIntf.GetByUserID(ctx, req)

	// migrator.Down()
	fmt.Println(res)
	assert.NotEmpty(t, res)
	assert.NoError(t, err)

}

// // success to buy merch
// func TestMerchBuySuccess_01(t *testing.T) {

// 	mockLogger := new(loggerMock.Interface)
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	cfg := NewTestConfig()

// 	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	userStrgIntf := postgres.NewUserStrg(dbConnector)
// 	boughtMerchStrgIntf := postgres.NewBoughtMerchStrg(dbConnector)

// 	authSvcIntf := service.NewAuthSvc(mockLogger, userStrgIntf, cfg.JwtKey)
// 	merchSvcIntf := service.NewMerchSvc(mockLogger, boughtMerchStrgIntf)

// 	req := &svcDto.SignInRequest{
// 		Username: "test5",
// 		Password: "test5",
// 	}

// 	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
// 	mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
// 	mockLogger.On("Warnf", mock.Anything, mock.Anything).Times(0)

// 	response, err := authSvcIntf.SignIn(ctx, req)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	token := response.JwtToken

// 	payload, err := jwtManager.VerifyAuthToken(token, cfg.Jwt.Key)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	id, err := uuid.Parse(payload.Id)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	err = merchSvcIntf.Buy(ctx, &svcDto.BuyMerchRequest{
// 		UserID:    id,
// 		MerchName: "t-shirt",
// 	})

// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, response)
// }

// // incorrect user id
// func TestMerchBuyFailed_01(t *testing.T) {

// 	mockLogger := new(loggerMock.Interface)
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	cfg := NewTestConfig()

// 	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	userStrgIntf := postgres.NewUserStrg(dbConnector)
// 	boughtMerchStrgIntf := postgres.NewBoughtMerchStrg(dbConnector)

// 	authSvcIntf := service.NewAuthSvc(mockLogger, userStrgIntf, cfg.JwtKey)
// 	merchSvcIntf := service.NewMerchSvc(mockLogger, boughtMerchStrgIntf)

// 	req := &svcDto.SignInRequest{
// 		Username: "test5",
// 		Password: "test5",
// 	}

// 	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
// 	mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
// 	mockLogger.On("Warnf", mock.Anything, mock.Anything).Times(0)

// 	response, err := authSvcIntf.SignIn(ctx, req)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	token := response.JwtToken

// 	_, err = jwt.VerifyAuthToken(token, cfg.JwtKey)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	id, err := uuid.NewRandom()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	err = merchSvcIntf.Buy(ctx, &svcDto.BuyMerchRequest{
// 		UserID:    id,
// 		MerchName: "t-shirt",
// 	})

// 	assert.Error(t, err)
// 	assert.Equal(t, postgres.Test_errBuyDecreaseCoinsAmount, err)
// }

// // incorrect merch type
// func TestMerchBuyFailed_02(t *testing.T) {

// 	mockLogger := new(loggerMock.Interface)
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	cfg := NewTestConfig()

// 	dbConnector, err := postgres.NewDbConn(ctx, &cfg.Database.Postgres)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	userStrgIntf := postgres.NewUserStrg(dbConnector)
// 	boughtMerchStrgIntf := postgres.NewBoughtMerchStrg(dbConnector)

// 	authSvcIntf := service.NewAuthSvc(mockLogger, userStrgIntf, cfg.JwtKey)
// 	merchSvcIntf := service.NewMerchSvc(mockLogger, boughtMerchStrgIntf)

// 	req := &svcDto.SignInRequest{
// 		Username: "test5",
// 		Password: "test5",
// 	}

// 	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
// 	mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
// 	mockLogger.On("Warnf", mock.Anything, mock.Anything).Times(0)

// 	response, err := authSvcIntf.SignIn(ctx, req)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	token := response.JwtToken

// 	payload, err := jwt.VerifyAuthToken(token, cfg.JwtKey)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	id, err := uuid.Parse(payload.Id)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	err = merchSvcIntf.Buy(ctx, &svcDto.BuyMerchRequest{
// 		UserID:    id,
// 		MerchName: "another avito merch",
// 	})

// 	assert.Error(t, err)
// 	assert.Equal(t, postgres.Test_errBuyDecreaseCoinsAmount, err)
// }
