package service

import (
	"collapp/module/user/model/web"
	"context"
)

type UserService interface {
	Create(ctx context.Context, request web.UserCreateRequest) web.UserResponse
	Update(ctx context.Context, request web.UserUpdateRequest) web.UserResponse
	Delete(ctx context.Context, userId int) web.UserResponse
	FindById(ctx context.Context, userId int) web.UserResponse
	FindAll(ctx context.Context) []web.UserResponse
	FindByEmail(ctx context.Context, userEmail string) web.UserLoginResponse
	UpdateToken(ctx context.Context, request web.UserTokenUpdateRequest) web.UserResponse
	Logout(ctx context.Context, userId int) web.UserResponse
}
