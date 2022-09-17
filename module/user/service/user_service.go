package service

import (
	"collapp/module/user/model"
	"context"
)

type UserService interface {
	Create(ctx context.Context, request model.UserCreateRequest) model.UserResponse
	Update(ctx context.Context, request model.UserUpdateRequest) model.UserResponse
	Delete(ctx context.Context, userId int) model.UserResponse
	SoftDelete(ctx context.Context, request model.UserDeleteRequest) model.UserResponse
	FindById(ctx context.Context, userId int) model.UserResponse
	FindAll(ctx context.Context) []model.UserResponse
	FindByEmail(ctx context.Context, userEmail string) model.UserLoginResponse
	FindByTokenRefresh(ctx context.Context, userTokenRefresh string) model.UserLoginResponse
	UpdateToken(ctx context.Context, request model.UserUpdateTokenRequest) model.UserResponse
	Logout(ctx context.Context, userId int) model.UserResponse
}
