package coin

import (
	"avito-task-2025/backend/internal/service"
	svcDto "avito-task-2025/backend/internal/service/dto"
	"avito-task-2025/backend/pkg/jwt"
	"avito-task-2025/backend/pkg/logger"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int32  `json:"amount"`
}

type SendCoinResponse struct {
}

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

	var httpReq SendCoinRequest
	var errorStr string

	err := ctx.ShouldBindJSON(&httpReq)
	if err != nil {
		errorStr = "incorrect request body"
		c.loggerIntf.Errorf("%s: %s", errorStr, err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errorStr})
		return
	}

	id, err := c.jwtMngIntf.GetStringClaimFromJWT(ctx.Request.Context(), "id")
	if err != nil {
		errorStr = "cant claim id from token"
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
		UserID:         uuID,
		ToUsername: httpReq.ToUser,
		CoinsAmount:    httpReq.Amount,
	}
	err = c.coinSvcIntf.Send(ctx.Request.Context(), req)
	if err != nil {
		c.loggerIntf.Errorf("failed to send coins: %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	ctx.Status(http.StatusOK)

}
