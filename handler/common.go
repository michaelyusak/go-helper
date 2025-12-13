package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/michaelyusak/go-helper/apperror"
	"github.com/michaelyusak/go-helper/helper"
)

type Common struct {
	appHealthy *bool
}

func NewCommon(appHealthy *bool) *Common {
	return &Common{
		appHealthy: appHealthy,
	}
}

func (h *Common) Ping(ctx *gin.Context) {
	helper.ResponseOK(ctx, "ok")
}

func (h *Common) NoRoute(ctx *gin.Context) {
	ctx.Error(apperror.NotFoundError())
}

func (h *Common) Health(ctx *gin.Context) {
	if !*h.appHealthy {
		ctx.Error(apperror.UnavailableError())
		return
	}

	helper.ResponseOK(ctx, "ok")
}
