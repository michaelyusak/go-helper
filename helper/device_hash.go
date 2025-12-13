package helper

import "fmt"

func GenerateDeviceHash(ipAddress, userAgent, deviceInfo string) string {
	return HashSHA512(fmt.Sprintf("%s:%s:%s", ipAddress, userAgent, deviceInfo))
}
