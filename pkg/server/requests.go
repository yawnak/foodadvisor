package server

type requestLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type requestJson struct {
	Id             int32  `json:"id"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	ExpirationDays int32  `json:"expiration"` //in days
	Role           string `json:"role"`
}
