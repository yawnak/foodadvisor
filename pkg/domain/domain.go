package domain

type User struct {
	Id             int64
	Username       string
	Password       string
	ExpirationDays int32 //in days
}
