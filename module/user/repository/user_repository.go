package repository

import (
	"collapp/module/user/model/domain"
	"context"
	"database/sql"
)

type UserRepository interface {
	Save(ctx context.Context, tx *sql.Tx, user domain.User) domain.User
	Update(ctx context.Context, tx *sql.Tx, user domain.User) domain.User
	Delete(ctx context.Context, tx *sql.Tx, user domain.User)
	SoftDelete(ctx context.Context, tx *sql.Tx, user domain.User)
	FindById(ctx context.Context, tx *sql.Tx, userId int) (domain.User, error)
	FindAll(ctx context.Context, tx *sql.Tx) []domain.User
	FindByEmail(ctx context.Context, tx *sql.Tx, userEmail string) (domain.User, error)
	FindByTokenRefresh(ctx context.Context, tx *sql.Tx, userTokenRefresh string) (domain.User, error)
	UpdateToken(ctx context.Context, tx *sql.Tx, user domain.User) domain.User
	Logout(ctx context.Context, tx *sql.Tx, user domain.User) domain.User
}
