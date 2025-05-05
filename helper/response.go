package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/michaelyusak/go-helper/appconstant"
	"github.com/michaelyusak/go-helper/dto"
)

func ResponseOK(ctx *gin.Context, res any) {
	ctx.JSON(http.StatusOK, dto.Response{Message: appconstant.MsgOK, Data: res})
}

func HealthOK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dto.GetHealthResponse{Healthy: true, Message: appconstant.MsgOK})
}
