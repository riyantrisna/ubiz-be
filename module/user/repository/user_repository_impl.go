package repository

import (
	"collapp/helper"
	"collapp/module/user/model/domain"
	"context"
	"database/sql"
	"errors"
)

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, user domain.User) domain.User {

	SQL := "INSERT INTO user(user_name, user_email, user_password) VALUES (?, ?, ?)"
	result, err := tx.ExecContext(ctx, SQL, user.UserName, user.UserEmail, user.UserPassword)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	user.UserId = int(id)
	return user
}

func (repository *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user domain.User) domain.User {
	SQL := "UPDATE user SET user_name = ?, user_email = ? WHERE user_id = ?"
	_, err := tx.ExecContext(ctx, SQL, user.UserName, user.UserEmail, user.UserId)
	helper.PanicIfError(err)

	return user
}

func (repository *UserRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, user domain.User) {
	SQL := "DELETE FROM user WHERE user_id = ?"
	_, err := tx.ExecContext(ctx, SQL, user.UserId)
	helper.PanicIfError(err)
}

func (repository *UserRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, userId int) (domain.User, error) {
	SQL := "SELECT user_id, user_name, user_email, IFNULL(user_token,''), IFNULL(user_token_refresh,'') FROM user WHERE user_id = ?"
	rows, err := tx.QueryContext(ctx, SQL, userId)
	helper.PanicIfError(err)
	defer rows.Close()

	user := domain.User{}
	if rows.Next() {
		err := rows.Scan(&user.UserId, &user.UserName, &user.UserEmail, &user.UserToken, &user.UserTokenRefresh)
		helper.PanicIfError(err)
		return user, nil
	} else {
		return user, errors.New("user is not found")
	}
}

func (repository *UserRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.User {
	SQL := "SELECT user_id, user_name, user_email, IFNULL(user_token,''), IFNULL(user_token_refresh,'') FROM user"
	rows, err := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(err)
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		user := domain.User{}
		err := rows.Scan(&user.UserId, &user.UserName, &user.UserEmail, &user.UserToken, &user.UserTokenRefresh)
		helper.PanicIfError(err)
		users = append(users, user)
	}

	return users
}

func (repository *UserRepositoryImpl) FindByEmail(ctx context.Context, tx *sql.Tx, userEmail string) (domain.User, error) {
	SQL := "SELECT user_id, user_name, user_password FROM user WHERE user_email = ?"
	rows, err := tx.QueryContext(ctx, SQL, userEmail)
	helper.PanicIfError(err)
	defer rows.Close()

	user := domain.User{}
	if rows.Next() {
		err := rows.Scan(&user.UserId, &user.UserName, &user.UserPassword)
		helper.PanicIfError(err)
		return user, nil
	} else {
		return user, errors.New("user is not found")
	}
}

func (repository *UserRepositoryImpl) UpdateToken(ctx context.Context, tx *sql.Tx, user domain.User) domain.User {
	SQL := "UPDATE user SET user_token = ?, user_token_refresh = ? WHERE user_id = ?"
	_, err := tx.ExecContext(ctx, SQL, user.UserToken, user.UserTokenRefresh, user.UserId)
	helper.PanicIfError(err)

	return user
}
