package entity

type SmtpHelperConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Identity string `json:"identity"`
}
