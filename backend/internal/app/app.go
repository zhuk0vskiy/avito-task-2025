package app

import (
	"avito-task-2025/backend/config"
	"avito-task-2025/backend/internal/service"
	"avito-task-2025/backend/internal/storage/postgres"
	"avito-task-2025/backend/pkg/jwt"
	"avito-task-2025/backend/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Logger       logger.Interface
	AuthSvcIntf  service.AuthIntf
	UserSvcIntf  service.UserIntf
	CoinSvcIntf  service.CoinIntf
	MerchSvcIntf service.MerchIntf
	JwtMngIntf   jwt.ManagerIntf
}

func NewApp(
	cfg *config.Config,
	logger logger.Interface,
	dbConn *pgxpool.Pool,
) *App {
	boughtMerchStrgIntf := postgres.NewBoughtMerchStrg(dbConn)
	transactionStrgIntf := postgres.NewTransactionStrg(dbConn)
	userStrgIntf := postgres.NewUserStrg(dbConn)

	jwtMngIntf := jwt.NewJwtManager(cfg.Jwt.Key, cfg.Jwt.ExpTimeHour)

	authSvcIntf := service.NewAuthSvc(logger, jwtMngIntf, userStrgIntf)
	userSvcIntf := service.NewUserSvc(logger, userStrgIntf, boughtMerchStrgIntf, transactionStrgIntf)
	coinSvcIntf := service.NewCoinSvc(logger, transactionStrgIntf)
	merchSvcIntf := service.NewMerchSvc(logger, boughtMerchStrgIntf)

	return &App{
		Logger:       logger,
		AuthSvcIntf:  authSvcIntf,
		UserSvcIntf:  userSvcIntf,
		CoinSvcIntf:  coinSvcIntf,
		MerchSvcIntf: merchSvcIntf,
		JwtMngIntf:      jwtMngIntf,
	}
}
