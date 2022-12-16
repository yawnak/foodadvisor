package domain

type User struct {
	Id             int32
	Username       string
	Password       string
	ExpirationDays int32 //in days
}
