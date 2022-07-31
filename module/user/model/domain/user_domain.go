package domain

type User struct {
	UserId           int
	UserName         string
	UserEmail        string
	UserPassword     string
	UserToken        string
	UserTokenRefresh string
	UserLangCode     string
}
