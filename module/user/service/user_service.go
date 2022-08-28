package service

import (
	"collapp/module/user/model/domain"
	"collapp/module/user/model/web"
	"context"
)

type UserService interface {
	Create(ctx context.Context, request web.UserCreateRequest) web.UserResponse
	Update(ctx context.Context, request web.UserUpdateRequest) web.UserResponse
	Delete(ctx context.Context, userId int) web.UserResponse
	SoftDelete(ctx context.Context, request web.UserDeleteRequest) web.UserResponse
	FindById(ctx context.Context, userId int) web.UserResponse
	FindAll(ctx context.Context) []web.UserResponse
	FindByEmail(ctx context.Context, userEmail string) web.UserLoginResponse
	FindByTokenRefresh(ctx context.Context, userTokenRefresh string) web.UserLoginResponse
	UpdateToken(ctx context.Context, request domain.User) web.UserResponse
	Logout(ctx context.Context, userId int) web.UserResponse
}
