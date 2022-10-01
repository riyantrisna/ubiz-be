package model

import (
	"database/sql"
	"mime/multipart"
)

// model User
type User struct {
	UserId                int
	UserName              string
	UserEmail             string
	UserPassword          string
	UserToken             string
	UserTokenCheck        sql.NullString
	UserTokenRefresh      string
	UserTokenRefreshCheck sql.NullString
	UserLangCode          string
	UserLastLogin         string
	UserLastLoginCheck    sql.NullString
	UserPhoto             string
	UserPhotoCheck        sql.NullString
	CreatedBy             int
	CreatedByCheck        sql.NullInt32
	CreatedByName         string
	CreatedByNameCheck    sql.NullString
	CreatedAt             string
	CreatedAtCheck        sql.NullString
	UpdatedBy             int
	UpdatedByCheck        sql.NullInt32
	UpdatedByName         string
	UpdatedByNameCheck    sql.NullString
	UpdatedAt             string
	UpdatedAtCheck        sql.NullString
	DeletedBy             int
	DeletedByName         string
	DeletedAt             string
}

// request
type UserCreateRequest struct {
	UserName      string                `validate:"required,min=1,max=200" form:"user_name"`
	UserEmail     string                `validate:"required,min=1,max=200,email" form:"user_email"`
	UserPassword  string                `validate:"required,min=1"`
	UserLangCode  string                `validate:"required,min=1" form:"user_lang_code"`
	UserPhoto     *multipart.FileHeader `form:"user_photo"`
	UserPhotoName string                `form:"-"`
	CreatedBy     int                   `validate:"required"`
	CreatedAt     string                `validate:"required"`
}

type UserUpdateRequest struct {
	UserId        int                   `validate:"required"`
	UserName      string                `validate:"required,min=1,max=200" form:"user_name"`
	UserEmail     string                `validate:"required,min=1,max=200,email" form:"user_email"`
	UserLangCode  string                `validate:"required,min=1" form:"user_lang_code"`
	UserPhoto     *multipart.FileHeader `form:"user_photo"`
	UserPhotoName string                `form:"-"`
	UpdatedBy     int                   `validate:"required"`
	UpdatedAt     string                `validate:"required"`
}

type UserDeleteRequest struct {
	UserId       int    `validate:"required"`
	DeletedBy    int    `validate:"required"`
	DeletedAt    string `validate:"required"`
	IsSoftDelete bool   `validate:"required" json:"is_soft_delete"`
}

type UserLoginRequest struct {
	UserEmail    string `validate:"required,min=1,email" json:"email"`
	UserPassword string `validate:"required,min=1" json:"password"`
}

type UserUpdateTokenRequest struct {
	UserId           int    `validate:"required"`
	UserToken        string `validate:"required"`
	UserTokenRefresh string `validate:"required"`
	UserLastLogin    string `validate:"required"`
}

// rersponse
type UserLoginResponse struct {
	UserId       int    `json:"user_id"`
	UserName     int    `json:"user_name"`
	UserPassword string `json:"user_password"`
	UserLangCode string `json:"user_lang_code"`
}

type UserResponse struct {
	UserId           int    `json:"user_id"`
	UserName         string `json:"user_name"`
	UserEmail        string `json:"user_email"`
	UserToken        string `json:"user_token"`
	UserTokenRefresh string `json:"user_token_refresh"`
	UserLangCode     string `json:"user_lang_code"`
	UserLastLogin    string `json:"user_last_login"`
	UserPhoto        string `json:"user_photo"`
	CreatedBy        int    `json:"created_by"`
	CreatedByName    string `json:"created_by_name"`
	CreatedAt        string `json:"created_at"`
	UpdatedBy        int    `json:"updated_by"`
	UpdatedByName    string `json:"updated_by_name"`
	UpdatedAt        string `json:"updated_at"`
}

func ToUserResponse(user User) UserResponse {
	return UserResponse{
		UserId:           user.UserId,
		UserName:         user.UserName,
		UserEmail:        user.UserEmail,
		UserToken:        user.UserToken,
		UserTokenRefresh: user.UserTokenRefresh,
		UserLangCode:     user.UserLangCode,
		UserLastLogin:    user.UserLastLogin,
		UserPhoto:        user.UserPhoto,
		CreatedBy:        user.CreatedBy,
		CreatedByName:    user.CreatedByName,
		CreatedAt:        user.CreatedAt,
		UpdatedBy:        user.UpdatedBy,
		UpdatedByName:    user.UpdatedByName,
		UpdatedAt:        user.UpdatedAt,
	}
}

func ToUserResponses(users []User) []UserResponse {
	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, ToUserResponse(user))
	}
	return userResponses
}

func ToUserLoginResponse(user User) UserLoginResponse {
	return UserLoginResponse{
		UserId:       user.UserId,
		UserPassword: user.UserPassword,
	}
}
