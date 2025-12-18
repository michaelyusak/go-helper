package middleware

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/michaelyusak/go-helper/appconstant"
	"github.com/michaelyusak/go-helper/apperror"
	"github.com/michaelyusak/go-helper/dto"
	"github.com/michaelyusak/go-helper/helper"
	"github.com/michaelyusak/go-helper/rest"
	"github.com/sirupsen/logrus"
)

type AuthOpt struct {
	IsCheckDeviceId    bool
	IsCheckAccessToken bool

	AllowedDeviceInfo []string
	AllowedIpAddress  []string

	AuthEngineBaseUrl string
}

type auth struct {
	isCheckDeviceId    bool
	isCheckAccessToken bool

	allowedDeviceInfo []string
	allowedIpAddress  []string

	authEngineRestClient rest.AuthRepo
}

func NewAuth(opt AuthOpt) *auth {
	return &auth{
		isCheckDeviceId:    opt.IsCheckDeviceId,
		isCheckAccessToken: opt.IsCheckAccessToken,

		allowedDeviceInfo: opt.AllowedDeviceInfo,
		allowedIpAddress:  opt.AllowedIpAddress,

		authEngineRestClient: rest.NewGoAuthRepo(rest.GoAuthRepoOpt{
			BaseUrl: opt.AuthEngineBaseUrl,
		}),
	}
}

func (m *auth) checkDeviceId(c *gin.Context, ipAddress, userAgent, deviceInfo, xDeviceId string) bool {
	if ipAddress == "" || userAgent == "" || deviceInfo == "" || xDeviceId == "" {
		return false
	}

	if len(m.allowedIpAddress) > 0 && !slices.Contains(m.allowedIpAddress, ipAddress) {
		return false
	}

	if len(m.allowedDeviceInfo) > 0 && !slices.Contains(m.allowedDeviceInfo, deviceInfo) {
		return false
	}

	ctx := helper.InjectValues(c.Request.Context(), map[appconstant.ContextKey]any{
		appconstant.UserAgentKey:      userAgent,
		appconstant.IpAddressKey:      ipAddress,
		appconstant.DeviceInfokey:     deviceInfo,
		appconstant.UniqueDeviceIdKey: xDeviceId,
	})

	c.Request = c.Request.WithContext(ctx)

	return true
}

func (m *auth) checkToken(c *gin.Context, isCheckToken bool) *apperror.AppError {
	authorization := c.Request.Header.Get(appconstant.Authorization)
	if authorization == "" && isCheckToken {
		return apperror.UnauthorizedError(apperror.AppErrorOpt{})
	}

	var token string

	if authorization != "" {
		split := strings.Split(authorization, " ")
		if len(split) != 2 || split[0] != appconstant.Bearer {
			return apperror.UnauthorizedError(apperror.AppErrorOpt{})
		}

		token = split[1]
	}

	ctx := helper.InjectValues(c.Request.Context(), map[appconstant.ContextKey]any{
		appconstant.AccessTokenKey: token,
	})
	c.Request = c.Request.WithContext(ctx)

	if !isCheckToken {
		return nil
	}

	customClaims, statusCode, err := m.authEngineRestClient.ValidateToken(c.Request.Context(), token)
	if err != nil {
		logrus.WithError(err).Error("[authMiddleware][checkToken][authEngineRestClient.ValidateToken] Error")
		return apperror.NewAppError(apperror.AppErrorOpt{
			Code:            statusCode,
			ResponseMessage: http.StatusText(statusCode),
			Message:         fmt.Sprintf("[authMiddleware][checkToken][authEngineRestClient.ValidateToken] Error: %v", err),
		})
	}

	ctx = helper.InjectValues(c.Request.Context(), map[appconstant.ContextKey]any{
		appconstant.AccountIdKey: customClaims.AccountId,
		appconstant.DeviceIdKey:  customClaims.DeviceId,
		appconstant.EmailKey:     customClaims.Email,
		appconstant.NameKey:      customClaims.Name,
	})
	c.Request = c.Request.WithContext(ctx)

	return nil
}

func (m *auth) Auth() func(c *gin.Context) {
	return func(c *gin.Context) {
		ipAddress := strings.TrimSpace(c.Request.Header.Get(appconstant.CfConnectingIp))
		if ipAddress == "" {
			ipAddress = c.ClientIP()
		}

		if m.isCheckDeviceId {
			passed := m.checkDeviceId(c, ipAddress, c.Request.UserAgent(), c.Request.Header.Get(appconstant.DeviceInfo), c.Request.Header.Get(appconstant.XDeviceId))
			if !passed {
				c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Message: appconstant.MsgUnauthorized})
				return
			}
		}

		err := m.checkToken(c, m.isCheckAccessToken)
		if err != nil {
			c.AbortWithStatusJSON(err.Code, dto.ErrorResponse{Message: err.ResponseMessage})
			return
		}

		c.Next()
	}
}
