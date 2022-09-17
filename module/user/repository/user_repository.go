package repository

import (
	"collapp/module/user/model"
	"context"
	"database/sql"
)

type UserRepository interface {
	Save(ctx context.Context, tx *sql.Tx, user model.UserCreateRequest) model.User
	Update(ctx context.Context, tx *sql.Tx, user model.UserUpdateRequest) model.User
	Delete(ctx context.Context, tx *sql.Tx, user model.User)
	SoftDelete(ctx context.Context, tx *sql.Tx, user model.User)
	FindById(ctx context.Context, tx *sql.Tx, userId int) (model.User, error)
	FindAll(ctx context.Context, tx *sql.Tx) []model.User
	FindByEmail(ctx context.Context, tx *sql.Tx, userEmail string) (model.User, error)
	FindByTokenRefresh(ctx context.Context, tx *sql.Tx, userTokenRefresh string) (model.User, error)
	UpdateToken(ctx context.Context, tx *sql.Tx, user model.User) model.User
	Logout(ctx context.Context, tx *sql.Tx, user model.User) model.User
}
