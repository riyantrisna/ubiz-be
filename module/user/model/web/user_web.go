package web

import "collapp/module/user/model/domain"

// request
type UserCreateRequest struct {
	UserName     string `validate:"required,min=1,max=200" json:"user_name"`
	UserEmail    string `validate:"required,min=1,max=200,email" json:"user_email"`
	UserPassword string `validate:"required,min=1"`
	UserLangCode string `validate:"required,min=1" json:"user_lang_code"`
}

type UserLoginRequest struct {
	UserEmail    string `validate:"required,min=1,email" json:"email"`
	UserPassword string `validate:"required,min=1" json:"password"`
}

type UserTokenUpdateRequest struct {
	UserId           int    `validate:"required"`
	UserToken        string `validate:"required" json:"user_token"`
	UserTokenRefresh string `validate:"required" json:"user_token_refresh"`
}

type UserUpdateRequest struct {
	UserId       int    `validate:"required"`
	UserName     string `validate:"required,min=1,max=200" json:"user_name"`
	UserEmail    string `validate:"required,min=1,max=200,email" json:"user_email"`
	UserLangCode string `validate:"required,min=1" json:"user_lang_code"`
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
}

func ToUserResponse(user domain.User) UserResponse {
	return UserResponse{
		UserId:           user.UserId,
		UserName:         user.UserName,
		UserEmail:        user.UserEmail,
		UserToken:        user.UserToken,
		UserTokenRefresh: user.UserTokenRefresh,
		UserLangCode:     user.UserLangCode,
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
