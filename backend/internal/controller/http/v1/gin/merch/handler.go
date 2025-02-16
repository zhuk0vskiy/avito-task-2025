package merch

import (
	"avito-task-2025/backend/internal/service"
	svcDto "avito-task-2025/backend/internal/service/dto"
	"avito-task-2025/backend/pkg/jwt"
	"avito-task-2025/backend/pkg/logger"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Controller struct {
	loggerIntf   logger.Interface
	merchSvcIntf service.MerchIntf
	jwtMngIntf   jwt.ManagerIntf
}

func NewMerchController(loggerIntf logger.Interface, merchSvcIntf service.MerchIntf, jwtMngIntf jwt.ManagerIntf) *Controller {
	return &Controller{
		loggerIntf:   loggerIntf,
		merchSvcIntf: merchSvcIntf,
		jwtMngIntf:   jwtMngIntf,
	}
}

func (c *Controller) BuyMerchHandler(ctx *gin.Context) {
	var errorStr string

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

	item := ctx.Param("item")

	req := &svcDto.BuyMerchRequest{
		UserID:    uuID,
		MerchName: item,
	}
	err = c.merchSvcIntf.Buy(ctx.Request.Context(), req)
	if err != nil {
		c.loggerIntf.Warnf("failed to buy merch: %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
