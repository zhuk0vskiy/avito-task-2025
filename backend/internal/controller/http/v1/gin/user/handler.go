package user

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
	userSvcIntf service.UserIntf
	jwtMngIntf  jwt.ManagerIntf
}

func NewUserController(loggerIntf logger.Interface, userSvcIntf service.UserIntf, jwtMngIntf jwt.ManagerIntf) *Controller {
	return &Controller{
		loggerIntf:  loggerIntf,
		userSvcIntf: userSvcIntf,
		jwtMngIntf:  jwtMngIntf,
	}
}

func (c *Controller) GetInfoHandler(ctx *gin.Context) {

	var errorStr string
	// fmt.Println(ctx.GetHeader("Authorization"))
	// fmt.Println(ctx.Request.Context())
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
	req := &svcDto.GetUserInfoRequest{
		UserID: uuID,
	}

	svcResponse, err := c.userSvcIntf.GetInfo(ctx.Request.Context(), req)
	if err != nil {
		c.loggerIntf.Warnf("failed to get user info: %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	inventory := make([]*v1.Inventory, 0)
	for i := 0; i < len(svcResponse.Inventory); i++ {
		inventory = append(inventory, &v1.Inventory{
			Type:     svcResponse.Inventory[i].Type,
			Quantity: svcResponse.Inventory[i].Amount,
		})
	}

	receive := make([]*v1.CoinHistoryReceived, 0)
	for i := 0; i < len(svcResponse.CoinHistory.Received); i++ {
		receive = append(receive, &v1.CoinHistoryReceived{
			FromUser: svcResponse.CoinHistory.Received[i].FromUsername,
			Amount:   svcResponse.CoinHistory.Received[i].CoinsAmount,
		})
	}

	sent := make([]*v1.CoinHistorySent, 0)
	for i := 0; i < len(svcResponse.CoinHistory.Sent); i++ {
		sent = append(sent, &v1.CoinHistorySent{
			ToUser: svcResponse.CoinHistory.Sent[i].ToUsername,
			Amount: svcResponse.CoinHistory.Sent[i].CoinsAmount,
		})
	}

	ctx.JSON(http.StatusOK, &v1.GetInfoResponse{
		Coins:     svcResponse.Coins,
		Inventory: inventory,
		CoinHistory: &v1.CoinHistory{
			Received: receive,
			Sent:     sent,
		},
	})
}
