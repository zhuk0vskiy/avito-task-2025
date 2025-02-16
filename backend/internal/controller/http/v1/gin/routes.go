package gin

import (
	"avito-task-2025/backend/internal/controller/http/v1/gin/auth"
	"avito-task-2025/backend/internal/controller/http/v1/gin/coin"
	"avito-task-2025/backend/internal/controller/http/v1/gin/merch"
	"avito-task-2025/backend/internal/controller/http/v1/gin/user"
	"avito-task-2025/backend/internal/service"
	"avito-task-2025/backend/pkg/jwt"
	"avito-task-2025/backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

func setAuthRoute(handler *gin.RouterGroup, l logger.Interface, authSvcIntf service.AuthIntf) {
	a := auth.NewAuthController(l, authSvcIntf)

	handler.POST("/auth", a.SignInHandler)
}

func setCoinRoute(handler *gin.RouterGroup, l logger.Interface, coinSvcIntf service.CoinIntf, jwtMngIntf jwt.ManagerIntf) {
	a := coin.NewCoinController(l, coinSvcIntf, jwtMngIntf)

	handler.POST("/sendCoin", a.SendCoinHandler)
}

func setMerchRoute(handler *gin.RouterGroup, l logger.Interface, merchSvcIntf service.MerchIntf, jwtMngIntf jwt.ManagerIntf) {
	a := merch.NewMerchController(l, merchSvcIntf, jwtMngIntf)

	handler.GET("/buy/:item", a.BuyMerchHandler)
}

func setUserRoute(handler *gin.RouterGroup, l logger.Interface, userSvcIntf service.UserIntf, jwtMngIntf jwt.ManagerIntf) {
	a := user.NewUserController(l, userSvcIntf, jwtMngIntf)

	handler.GET("/info", a.GetInfoHandler)
}

func SetRoutes(
	handler *gin.RouterGroup,
	logger logger.Interface,
	authSvcIntf service.AuthIntf,
	coinSvcIntf service.CoinIntf,
	merchSvcIntf service.MerchIntf,
	userSvcIntf service.UserIntf,
	jwtMngIntf jwt.ManagerIntf,
) {
	setAuthRoute(handler, logger, authSvcIntf)
	setCoinRoute(handler, logger, coinSvcIntf, jwtMngIntf)
	setMerchRoute(handler, logger, merchSvcIntf, jwtMngIntf)
	setUserRoute(handler, logger, userSvcIntf, jwtMngIntf)

}
