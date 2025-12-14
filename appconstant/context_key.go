package appconstant

type ContextKey string

const (
	RequestId = "request-id"

	// Auth
	Bearer                     = "Bearer"
	DeviceIdKey     ContextKey = "device_id"
	AccessTokenKey  ContextKey = "access_token"
	RefreshTokenKey ContextKey = "refresh_token"
	AccountIdKey    ContextKey = "account_id"
	DeviceHashKey   ContextKey = "device_hash_key"
	EmailKey        ContextKey = "email"
	NameKey         ContextKey = "name"
	UserAgentKey    ContextKey = "user_agent"
	IpAddressKey    ContextKey = "ip_adress"
	DeviceInfokey   ContextKey = "device_info"

	// HTTP Header Key
	Authorization  = "Authorization"
	CfConnectingIp = "CF-Connecting-IP"
	DeviceInfo     = "Device-Info"
	UserAgent      = "User-Agent"
	ClientIp       = "X-Client-IP"
)
