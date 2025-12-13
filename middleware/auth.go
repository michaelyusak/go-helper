package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/michaelyusak/go-helper/appconstant"
	"github.com/michaelyusak/go-helper/dto"
	"github.com/michaelyusak/go-helper/helper"
)

type AuthOpt struct {
	IsCheckDeviceId bool

	Hash helper.HashHelper
}

type auth struct {
	isCheckDeviceId bool

	hash helper.HashHelper
}

func NewAuth(opt AuthOpt) *auth {
	return &auth{
		isCheckDeviceId: opt.IsCheckDeviceId,
		hash:            opt.Hash,
	}
}

func (m *auth) checkDeviceId(c *gin.Context, ipAddress, userAgent, referer string) {
	if ipAddress == "" || userAgent == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Message: appconstant.MsgUnauthorized})
		return
	}

	deviceId := helper.HashSHA512(fmt.Sprintf("%s:%s:%s", ipAddress, userAgent, referer))

	c.Set(string(appconstant.DeviceIdKey), deviceId)

	ctx := context.WithValue(c.Request.Context(), appconstant.DeviceIdKey, deviceId)
	c.Request = c.Request.WithContext(ctx)
}

func (m *auth) Auth() func(c *gin.Context) {
	return func(c *gin.Context) {
		ipAddress := strings.TrimSpace(c.Request.Header.Get(appconstant.CfConnectingIp))
		if ipAddress == "" {
			ipAddress = c.ClientIP()
		}

		if m.isCheckDeviceId {
			m.checkDeviceId(c, ipAddress, c.Request.UserAgent(), c.Request.Referer())
		}

		c.Next()
	}
}
