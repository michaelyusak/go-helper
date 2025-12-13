package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/michaelyusak/go-helper/appconstant"
	"github.com/michaelyusak/go-helper/dto"
	"github.com/michaelyusak/go-helper/helper"
)

type AuthOpt struct {
	IsCheckDeviceId   bool
	AllowedDeviceInfo []string
	AllowedIpAddress  []string
}

type auth struct {
	isCheckDeviceId   bool
	allowedDeviceInfo []string
	allowedIpAddress  []string

	hash helper.HashHelper
}

func NewAuth(opt AuthOpt) *auth {
	return &auth{
		isCheckDeviceId:   opt.IsCheckDeviceId,
		allowedDeviceInfo: opt.AllowedDeviceInfo,
	}
}

func (m *auth) checkDeviceId(c *gin.Context, ipAddress, userAgent, deviceInfo string) bool {
	if ipAddress == "" || userAgent == "" || deviceInfo == "" {
		return false
	}

	if len(m.allowedIpAddress) > 0 || !slices.Contains(m.allowedIpAddress, ipAddress) {
		return false
	}

	if len(m.allowedDeviceInfo) > 0 || !slices.Contains(m.allowedDeviceInfo, deviceInfo) {
		return false
	}

	deviceHash := helper.GenerateDeviceHash(ipAddress, userAgent, deviceInfo)

	ctx := helper.InjectValues(c.Request.Context(), map[appconstant.ContextKey]any{
		appconstant.DeviceHashKey: deviceHash,
		appconstant.UserAgentKey:  userAgent,
		appconstant.IpAddressKey:  ipAddress,
		appconstant.DeviceInfokey: deviceInfo,
	})

	c.Request = c.Request.WithContext(ctx)

	return true
}

func (m *auth) Auth() func(c *gin.Context) {
	return func(c *gin.Context) {
		ipAddress := strings.TrimSpace(c.Request.Header.Get(appconstant.CfConnectingIp))
		if ipAddress == "" {
			ipAddress = c.ClientIP()
		}

		if m.isCheckDeviceId {
			passed := m.checkDeviceId(c, ipAddress, c.Request.UserAgent(), c.Request.Header.Get(appconstant.DeviceInfo))
			if !passed {
				c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Message: appconstant.MsgUnauthorized})
				return
			}
		}

		c.Next()
	}
}
