package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/michaelyusak/go-helper/apperror"
	"github.com/michaelyusak/go-helper/helper"
)

type CommonHandler struct {
	appHealthy *bool
}

func NewCommonHandler(appHealthy *bool) *CommonHandler {
	return &CommonHandler{
		appHealthy: appHealthy,
	}
}

func (h *CommonHandler) Ping(ctx *gin.Context) {
	helper.ResponseOK(ctx, "ok")
}

func (h *CommonHandler) NoRoute(ctx *gin.Context) {
	ctx.Error(apperror.NotFoundError())
}

func (h *CommonHandler) Health(ctx *gin.Context) {
	if !*h.appHealthy {
		ctx.Error(apperror.UnavailableError())
		return
	}

	helper.ResponseOK(ctx, "ok")
}
