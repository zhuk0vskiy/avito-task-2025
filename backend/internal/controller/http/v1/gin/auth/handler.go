package auth

import (
	v1 "avito-task-2025/backend/internal/controller/http/v1"
	svcDto "avito-task-2025/backend/internal/service/dto"
	"avito-task-2025/backend/pkg/logger"

	"avito-task-2025/backend/internal/service"

	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	loggerIntf  logger.Interface
	authSvcIntf service.AuthIntf
}

func NewAuthController(loggerIntf logger.Interface, authSvcIntf service.AuthIntf) *Controller {
	return &Controller{
		loggerIntf:  loggerIntf,
		authSvcIntf: authSvcIntf,
	}
}

func (c *Controller) SignInHandler(ctx *gin.Context) {

	//start := time.Now()

	// wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

	var httpReq v1.SignInRequest
	var errorStr string

	if err := ctx.ShouldBindJSON(&httpReq); err != nil {
		errorStr = "incorrect request body"
		c.loggerIntf.Warnf("%s: %s", errorStr, err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errorStr})
		return
	}

	req := &svcDto.SignInRequest{
		Username: httpReq.Username,
		Password: httpReq.Password,
	}
	svcResponse, err := c.authSvcIntf.SignIn(ctx.Request.Context(), req)
	if err != nil {
		c.loggerIntf.Warnf("cant auth: %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, &v1.SignInResponse{
		Token: svcResponse.JwtToken,
	})

}
