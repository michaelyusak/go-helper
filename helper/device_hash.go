package helper

import (
	"context"
	"fmt"

	"github.com/michaelyusak/go-helper/appconstant"
)

func GenerateDeviceHash(ctx context.Context) string {
	deviceInfo := ctx.Value(appconstant.DeviceInfokey).(string)
	uniqueDeviceId := ctx.Value(appconstant.UniqueDeviceIdKey).(string)
	accountId := ctx.Value(appconstant.AccountIdKey).(int64)

	return HashSHA512(fmt.Sprintf("%v:%s:%s", accountId, deviceInfo, uniqueDeviceId))
}
