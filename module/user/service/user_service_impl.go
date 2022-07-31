package service

import (
	"collapp/helper"
	"collapp/module/user/model/domain"
	"collapp/module/user/model/web"
	"collapp/module/user/repository"
	"context"
	"database/sql"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	DB             *sql.DB
}

func NewUserService(DB *sql.DB) UserService {
	userRepository := repository.NewUserRepository()

	return &UserServiceImpl{
		UserRepository: userRepository,
		DB:             DB,
	}
}

func (service *UserServiceImpl) Create(ctx context.Context, request web.UserCreateRequest) web.UserResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	userData := domain.User{
		UserName:     request.UserName,
		UserEmail:    request.UserEmail,
		UserPassword: request.UserPassword,
		UserLangCode: request.UserLangCode,
	}

	userData = service.UserRepository.Save(ctx, tx, userData)

	return web.ToUserResponse(userData)
}

func (service *UserServiceImpl) Update(ctx context.Context, request web.UserUpdateRequest) web.UserResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	userData, err := service.UserRepository.FindById(ctx, tx, request.UserId)
	if err == nil {
		userData.UserName = request.UserName
		userData.UserEmail = request.UserEmail
		userData.UserLangCode = request.UserLangCode

		userData = service.UserRepository.Update(ctx, tx, userData)
	}

	return web.ToUserResponse(userData)
}

func (service *UserServiceImpl) Delete(ctx context.Context, userId int) web.UserResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	userData, err := service.UserRepository.FindById(ctx, tx, userId)
	if err == nil {
		service.UserRepository.Delete(ctx, tx, userData)
	}

	return web.ToUserResponse(userData)
}

func (service *UserServiceImpl) FindById(ctx context.Context, userId int) web.UserResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	userData, _ := service.UserRepository.FindById(ctx, tx, userId)

	return web.ToUserResponse(userData)
}

func (service *UserServiceImpl) FindAll(ctx context.Context) []web.UserResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	usersData := service.UserRepository.FindAll(ctx, tx)

	return web.ToUserResponses(usersData)
}

func (service *UserServiceImpl) FindByEmail(ctx context.Context, userEmail string) web.UserLoginResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	userData, _ := service.UserRepository.FindByEmail(ctx, tx, userEmail)

	return web.ToUserLoginResponse(userData)
}

func (service *UserServiceImpl) UpdateToken(ctx context.Context, request web.UserTokenUpdateRequest) web.UserResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	userData, err := service.UserRepository.FindById(ctx, tx, request.UserId)
	if err == nil {
		userData.UserId = request.UserId
		userData.UserToken = request.UserToken
		userData.UserTokenRefresh = request.UserTokenRefresh

		userData = service.UserRepository.UpdateToken(ctx, tx, userData)
	}

	return web.ToUserResponse(userData)
}
