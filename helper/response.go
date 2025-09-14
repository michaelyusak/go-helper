package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/michaelyusak/go-helper/appconstant"
	"github.com/michaelyusak/go-helper/dto"
)

func ResponseOK(ctx *gin.Context, res any) {
	ctx.JSON(http.StatusOK, dto.Response{Message: appconstant.MsgOK, Success: true, StatusCode: http.StatusOK, Data: res})
}

func HealthOK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dto.GetHealthResponse{Health: appconstant.MsgOK})
}
