package appconstant

type ContextKey string

const (
	RequestId = "request-id"

	// Auth
	Bearer                 = "Bearer"
	DeviceIdKey ContextKey = "device_id"

	// HTTP Header Key
	Authorization  = "Authorization"
	CfConnectingIp = "CF-Connecting-IP"
)
