package controller

import (
	v1 "avito-task-2025/backend/internal/controller/http/v1"
	"avito-task-2025/backend/internal/service"
	"avito-task-2025/backend/pkg/jwt"
	"avito-task-2025/backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	handler      *gin.Engine
	routerGroups map[string]*gin.RouterGroup
}

func disableCors(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "*")
	ctx.Header("Access-Control-Allow-Headers", "*")
}

func NewRouter(handler *gin.Engine) *Controller {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Disable CORS
	handler.OPTIONS("/*any", disableCors)

	v1 := handler.Group("/api")

	return &Controller{
		handler: handler,
		routerGroups: map[string]*gin.RouterGroup{
			"v1": v1,
		},
	}
}

func (c *Controller) SetV1Routes(
	logger logger.Interface,
	authSvcIntf service.AuthIntf,
	coinSvcIntf service.CoinIntf,
	merchSvcIntf service.MerchIntf,
	userSvcIntf service.UserIntf,
	jwtMngIntf jwt.ManagerIntf,
) {
	v1.SetRoutes(
		c.routerGroups["v1"],
		logger,
		authSvcIntf,
		coinSvcIntf,
		merchSvcIntf,
		userSvcIntf,
		jwtMngIntf,
	)
}
