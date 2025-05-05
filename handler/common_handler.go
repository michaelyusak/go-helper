package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/michaelyusak/go-helper/apperror"
	"github.com/michaelyusak/go-helper/helper"
)

var IsServiceHealthy = false

type CommonHandler struct{}

func (h *CommonHandler) Ping(ctx *gin.Context) {
	helper.ResponseOK(ctx, "ok")
}

func (h *CommonHandler) NoRoute(ctx *gin.Context) {
	ctx.Error(apperror.NotFoundError())
}

func (h *CommonHandler) GetHealth(ctx *gin.Context) {
	if !IsServiceHealthy {
		ctx.Error(apperror.NewAppError(apperror.AppErrorOpt{
			Code:            http.StatusServiceUnavailable,
			ResponseMessage: "service unavailable",
			Message:         "service unavailable",
		}))
		return
	}

	helper.HealthOK(ctx)
}
