package server

type requestLogin struct {
	Username string `json:"username" validate:"required,alphanumwithunderscore,usernamemax"`
	Password string `json:"password" validate:"required,alphanumwithunderscore,passwordmax"`
}

type setUserRoleRequest struct {
	Role string `json:"role"`
}

type requestSignup struct {
	Id             int32   `json:"id" validate:"isdefault"`
	Username       string  `json:"username" validate:"required,alphanumwithunderscore,usernamemax"`
	Password       string  `json:"password" validate:"required,alphanumwithunderscore,passwordmax"`
	ExpirationDays int32   `json:"expiration" validate:"gte=0"` //in days
	Role           *string `json:"-" validate:"isdefault"`
}

type requestUpdateUser struct {
	Id             int32   `json:"-" validate:"isdefault"`
	Username       string  `json:"username" validate:"required,alphanumwithunderscore,usernamemax"`
	Password       string  `json:"-" validate:"isdefault"`
	ExpirationDays int32   `json:"expiration" validate:"gte=0"` //in days
	Role           *string `json:"-" validate:"isdefault"`
}


type requestCreateFood struct {
	Id   int32  `json:"id" validate:"isdefault"`
	Name string `json:"name"`
	//in minutes
	CookTime int32 `json:"cooktime" validate:"gt=0"`
}
