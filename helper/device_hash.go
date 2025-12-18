package helper

import (
	"context"
	"fmt"

	"github.com/michaelyusak/go-helper/appconstant"
)

func GenerateDeviceHash(ctx context.Context, accountId int64) string {
	deviceInfo := ctx.Value(appconstant.DeviceInfokey).(string)
	uniqueDeviceId := ctx.Value(appconstant.UniqueDeviceIdKey).(string)

	return HashSHA512(fmt.Sprintf("%v:%s:%s", accountId, deviceInfo, uniqueDeviceId))
}
