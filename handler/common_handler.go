package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/michaelyusak/go-helper/apperror"
	"github.com/michaelyusak/go-helper/helper"
)

type CommonHandler struct{}

func (h *CommonHandler) Ping(ctx *gin.Context) {
	helper.ResponseOK(ctx, "ok")
}

func (h *CommonHandler) NoRoute(ctx *gin.Context) {
	ctx.Error(apperror.NotFoundError())
}
