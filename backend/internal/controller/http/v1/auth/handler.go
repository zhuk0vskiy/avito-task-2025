package auth

import (
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

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

func (c *Controller) SignInHandler(ctx *gin.Context) {

	//start := time.Now()

	// wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

	var httpReq SignInRequest
	var errorStr string

	if err := ctx.ShouldBindJSON(&httpReq); err != nil {
		errorStr = "incorrect request body"
		c.loggerIntf.Errorf("%s: %s", errorStr, err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errorStr})
		return
	}

	req := &svcDto.SignInRequest{
		Username: httpReq.Username,
		Password: httpReq.Password,
	}
	svcResponse, err := c.authSvcIntf.SignIn(ctx.Request.Context(), req)
	if err != nil {
		c.loggerIntf.Errorf("cant auth: %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	ctx.JSON(http.StatusOK, &SignInResponse{
		Token: svcResponse.JwtToken,
	})

}
