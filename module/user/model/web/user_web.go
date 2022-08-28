package web

import (
	"collapp/module/user/model/domain"
)

// request
type UserCreateRequest struct {
	UserName     string `validate:"required,min=1,max=200" json:"user_name"`
	UserEmail    string `validate:"required,min=1,max=200,email" json:"user_email"`
	UserPassword string `validate:"required,min=1"`
	UserLangCode string `validate:"required,min=1" json:"user_lang_code"`
	CreatedBy    int    `validate:"required"`
	CreatedAt    string `validate:"required"`
}

type UserUpdateRequest struct {
	UserId       int    `validate:"required"`
	UserName     string `validate:"required,min=1,max=200" json:"user_name"`
	UserEmail    string `validate:"required,min=1,max=200,email" json:"user_email"`
	UserLangCode string `validate:"required,min=1" json:"user_lang_code"`
	UpdatedBy    int    `validate:"required"`
	UpdatedAt    string `validate:"required"`
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
	CreatedBy        int    `json:"created_by"`
	CreatedByName    string `json:"created_by_name"`
	CreatedAt        string `json:"created_at"`
	UpdatedBy        int    `json:"updated_by"`
	UpdatedByName    string `json:"updated_by_name"`
	UpdatedAt        string `json:"updated_at"`
	DeletedBy        int    `json:"deleted_by"`
	DeletedByName    string `json:"deleted_by_name"`
	DeletedAt        string `json:"deleted_at"`
}

func ToUserResponse(user domain.User) UserResponse {
	return UserResponse{
		UserId:           user.UserId,
		UserName:         user.UserName,
		UserEmail:        user.UserEmail,
		UserToken:        user.UserToken,
		UserTokenRefresh: user.UserTokenRefresh,
		UserLangCode:     user.UserLangCode,
		UserLastLogin:    user.UserLastLogin,
		CreatedBy:        user.CreatedBy,
		CreatedByName:    user.CreatedByName,
		CreatedAt:        user.CreatedAt,
		UpdatedBy:        user.UpdatedBy,
		UpdatedByName:    user.UpdatedByName,
		UpdatedAt:        user.UpdatedAt,
		DeletedBy:        user.DeletedBy,
		DeletedByName:    user.DeletedByName,
		DeletedAt:        user.DeletedAt,
	}
}

func ToUserResponses(users []domain.User) []UserResponse {
	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, ToUserResponse(user))
	}
	return userResponses
}

func ToUserLoginResponse(user domain.User) UserLoginResponse {
	return UserLoginResponse{
		UserId:       user.UserId,
		UserPassword: user.UserPassword,
	}
}
