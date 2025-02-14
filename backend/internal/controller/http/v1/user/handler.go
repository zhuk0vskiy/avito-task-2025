package user

import (
	"avito-task-2025/backend/internal/service"
	svcDto "avito-task-2025/backend/internal/service/dto"
	"avito-task-2025/backend/pkg/jwt"
	"avito-task-2025/backend/pkg/logger"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetInfoRequest struct {
}

type GetInfoInventory struct {
	Type     string `json:"type"`
	Quantity int16  `json:"quantity"`
}

type CoinHistoryReceived struct {
	FromUser string `json:"fromUser"`
	Amount   int32  `json:"amount"`
}

type CoinHistorySent struct {
	ToUser string `json:"toUser"`
	Amount int32  `json:"amount"`
}

type CoinHistory struct {
	Received []*CoinHistoryReceived `json:"received"`
	Sent     []*CoinHistorySent     `json:"sent"`
}

type GetInfoResponse struct {
	Coins       int32               `json:"coins"`
	Inventory   []*GetInfoInventory `json:"inventory"`
	CoinHistory *CoinHistory        `json:"coinHistory"`
}

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
	req := &svcDto.GetUserInfoRequest{
		UserID: uuID,
	}

	svcResponse, err := c.userSvcIntf.GetInfo(ctx.Request.Context(), req)
	if err != nil {
		c.loggerIntf.Errorf("failed to get user info: %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	inventory := make([]*GetInfoInventory, 0)
	for i := 0; i < len(svcResponse.Inventory); i++ {
		inventory = append(inventory, &GetInfoInventory{
			Type:     svcResponse.Inventory[i].Type,
			Quantity: svcResponse.Inventory[i].Amount,
		})
	}

	receive := make([]*CoinHistoryReceived, 0)
	for i := 0; i < len(svcResponse.CoinHistory.Received); i++ {
		receive = append(receive, &CoinHistoryReceived{
			FromUser: svcResponse.CoinHistory.Received[i].FromUsername,
			Amount:   svcResponse.CoinHistory.Received[i].CoinsAmount,
		})
	}

	sent := make([]*CoinHistorySent, 0)
	for i := 0; i < len(svcResponse.CoinHistory.Sent); i++ {
		sent = append(sent, &CoinHistorySent{
			ToUser: svcResponse.CoinHistory.Sent[i].ToUsername,
			Amount: svcResponse.CoinHistory.Sent[i].CoinsAmount,
		})
	}

	ctx.JSON(http.StatusOK, &GetInfoResponse{
		Coins:     svcResponse.Coins,
		Inventory: inventory,
		CoinHistory: &CoinHistory{
			Received: receive,
			Sent:     sent,
		},
	})
}
