package server

type requestLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
