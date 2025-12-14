package entity

type JwtCustomClaims struct {
	AccountId int64  `json:"account_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	DeviceId  int64  `json:"device_id"`
}
