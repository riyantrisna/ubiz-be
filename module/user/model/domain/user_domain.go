package domain

type User struct {
	UserId           int
	UserName         string
	UserEmail        string
	UserPassword     string
	UserToken        string
	UserTokenRefresh string
	UserLangCode     string
	UserLastLogin    string
	CreatedBy        int
	CreatedByName    string
	CreatedAt        string
	UpdatedBy        int
	UpdatedByName    string
	UpdatedAt        string
}
