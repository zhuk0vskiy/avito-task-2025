package coin

import (
	v1 "avito-task-2025/backend/internal/controller/http/v1"
	"avito-task-2025/backend/internal/service"
	svcDto "avito-task-2025/backend/internal/service/dto"
	"avito-task-2025/backend/pkg/jwt"
	"avito-task-2025/backend/pkg/logger"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Controller struct {
	loggerIntf  logger.Interface
	coinSvcIntf service.CoinIntf
	jwtMngIntf  jwt.ManagerIntf
}

func NewCoinController(loggerIntf logger.Interface, coinSvcIntf service.CoinIntf, jwtMngIntf jwt.ManagerIntf) *Controller {
	return &Controller{
		loggerIntf:  loggerIntf,
		coinSvcIntf: coinSvcIntf,
		jwtMngIntf:  jwtMngIntf,
	}
}

func (c *Controller) SendCoinHandler(ctx *gin.Context) {

	var httpReq v1.SendCoinRequest
	var errorStr string

	err := ctx.ShouldBindJSON(&httpReq)
	if err != nil {
		errorStr = "incorrect request body"
		c.loggerIntf.Warnf("%s: %s", errorStr, err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errorStr})
		return
	}

	id, err := c.jwtMngIntf.GinGetStringClaimFromJWT(ctx, "id")
	if err != nil {
		errorStr = "jwt token is invalid"
		c.loggerIntf.Errorf("%s: %s", errorStr, err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errorStr})
		return
	}

	uuID, err := uuid.Parse(id)
	if err != nil {
		errorStr = "failed to parse id"
		c.loggerIntf.Errorf("%s: %s", errorStr, err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errorStr})
		return
	}
	req := &svcDto.SendCoinsRequest{
		UserID:      uuID,
		ToUsername:  httpReq.ToUser,
		CoinsAmount: httpReq.Amount,
	}
	err = c.coinSvcIntf.Send(ctx.Request.Context(), req)
	if err != nil {
		c.loggerIntf.Warnf("failed to send coins: %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)

}
