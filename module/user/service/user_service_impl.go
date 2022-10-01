package service

import (
	"collapp/helper"
	"collapp/module/user/model"
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

func (service *UserServiceImpl) Create(ctx context.Context, request model.UserCreateRequest) model.UserResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	userData := service.UserRepository.Save(ctx, tx, request)
	if userData.UserId > 0 {
		userData, err := service.UserRepository.FindById(ctx, tx, userData.UserId)
		helper.PanicIfError(err)

		return model.ToUserResponse(userData)
	} else {
		return model.ToUserResponse(userData)
	}
}

func (service *UserServiceImpl) Update(ctx context.Context, request model.UserUpdateRequest) (model.UserResponse, string) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	userData, err := service.UserRepository.FindById(ctx, tx, request.UserId)
	userPhoto := userData.UserPhoto
	if err == nil {
		userData = service.UserRepository.Update(ctx, tx, request)

		userData, err := service.UserRepository.FindById(ctx, tx, userData.UserId)
		helper.PanicIfError(err)

		return model.ToUserResponse(userData), userPhoto
	}

	return model.ToUserResponse(userData), userPhoto
}

func (service *UserServiceImpl) Delete(ctx context.Context, userId int) model.UserResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	userData, err := service.UserRepository.FindById(ctx, tx, userId)
	if err == nil {
		service.UserRepository.Delete(ctx, tx, userData)
	}

	return model.ToUserResponse(userData)
}

func (service *UserServiceImpl) SoftDelete(ctx context.Context, request model.UserDeleteRequest) model.UserResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	userData, err := service.UserRepository.FindById(ctx, tx, request.UserId)
	if err == nil {
		userData.UserId = request.UserId
		userData.DeletedBy = request.DeletedBy
		userData.DeletedAt = request.DeletedAt

		service.UserRepository.SoftDelete(ctx, tx, userData)
	}

	return model.ToUserResponse(userData)
}

func (service *UserServiceImpl) FindById(ctx context.Context, userId int) model.UserResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	userData, _ := service.UserRepository.FindById(ctx, tx, userId)

	return model.ToUserResponse(userData)
}

func (service *UserServiceImpl) FindAll(ctx context.Context) []model.UserResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	usersData := service.UserRepository.FindAll(ctx, tx)

	return model.ToUserResponses(usersData)
}

func (service *UserServiceImpl) FindByEmail(ctx context.Context, userEmail string) model.UserLoginResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	userData, _ := service.UserRepository.FindByEmail(ctx, tx, userEmail)

	return model.ToUserLoginResponse(userData)
}

func (service *UserServiceImpl) FindByTokenRefresh(ctx context.Context, userTokenRefresh string) model.UserLoginResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	userData, _ := service.UserRepository.FindByTokenRefresh(ctx, tx, userTokenRefresh)

	return model.ToUserLoginResponse(userData)
}

func (service *UserServiceImpl) UpdateToken(ctx context.Context, request model.UserUpdateTokenRequest) model.UserResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	userData, err := service.UserRepository.FindById(ctx, tx, request.UserId)
	if err == nil {
		userData.UserId = request.UserId
		userData.UserToken = request.UserToken
		userData.UserTokenRefresh = request.UserTokenRefresh
		userData.UserLastLogin = request.UserLastLogin

		userData = service.UserRepository.UpdateToken(ctx, tx, userData)
	}

	return model.ToUserResponse(userData)
}

func (service *UserServiceImpl) Logout(ctx context.Context, userId int) model.UserResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	userData, err := service.UserRepository.FindById(ctx, tx, userId)
	if err == nil {
		service.UserRepository.Logout(ctx, tx, userData)
	}

	return model.ToUserResponse(userData)
}
