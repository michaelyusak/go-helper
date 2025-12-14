package helper

import (
	"context"

	"github.com/michaelyusak/go-helper/appconstant"
)

func InjectValues(ctx context.Context, values map[appconstant.ContextKey]any) context.Context {
	for k, v := range values {
		ctx = context.WithValue(ctx, k, v)
	}

	return ctx
}

func AuthHeadersFromContext(ctx context.Context) map[string]string {
	return map[string]string{
		appconstant.ClientIp:   ctx.Value(appconstant.IpAddressKey).(string),
		appconstant.DeviceInfo: ctx.Value(appconstant.DeviceInfo).(string),
		appconstant.UserAgent:  ctx.Value(appconstant.UserAgentKey).(string),
	}
}
