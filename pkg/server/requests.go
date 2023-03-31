package server

type requestLogin struct {
	Username string `json:"username" validate:"required,alphanumwithunderscore,usernamemax"`
	Password string `json:"password" validate:"required,alphanumwithunderscore,passwordmax"`
}
}

type requestSignup struct {
	Id             int32   `json:"id"`
	Username       string  `json:"username" validate:"required,alphanumwithunderscore,usernamemax"`
	Password       string  `json:"password" validate:"required,alphanumwithunderscore,passwordmax"`
	ExpirationDays int32   `json:"expiration" validate:"gte=0"` //in days
	Role           *string `json:"-" validate:"isdefault"`
}
