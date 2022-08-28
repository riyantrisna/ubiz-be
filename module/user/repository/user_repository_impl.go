package repository

import (
	"collapp/helper"
	"collapp/module/user/model/domain"
	"context"
	"database/sql"
)

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, user domain.User) domain.User {

	SQL := `INSERT INTO user
			(
				user_name, 
				user_email, 
				user_password, 
				user_lang_code, 
				created_by, 
				created_at
			) VALUES (
				?, 
				?, 
				?, 
				?, 
				?, 
				?
			)`
	result, err := tx.ExecContext(ctx, SQL,
		user.UserName,
		user.UserEmail,
		user.UserPassword,
		user.UserLangCode,
		user.CreatedBy,
		user.CreatedAt)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	user.UserId = int(id)
	return user
}

func (repository *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user domain.User) domain.User {
	SQL := `UPDATE 
				user 
			SET 
				user_name = ?, 
				user_email = ?, 
				user_lang_code = ?, 
				updated_by = ?, 
				updated_at = ? 
			WHERE 
				user_id = ?`
	_, err := tx.ExecContext(ctx, SQL,
		user.UserName,
		user.UserEmail,
		user.UserLangCode,
		user.UpdatedBy,
		user.UpdatedAt,
		user.UserId)
	helper.PanicIfError(err)

	return user
}

func (repository *UserRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, user domain.User) {
	SQL := `DELETE FROM user WHERE user_id = ?`
	_, err := tx.ExecContext(ctx, SQL, user.UserId)
	helper.PanicIfError(err)
}

func (repository *UserRepositoryImpl) SoftDelete(ctx context.Context, tx *sql.Tx, user domain.User) {
	SQL := `UPDATE 
				user 
			SET 
				deleted_by = ?, 
				deleted_at = ? 
			WHERE 
				user_id = ?`
	_, err := tx.ExecContext(ctx, SQL,
		user.DeletedBy,
		user.DeletedAt,
		user.UserId)
	helper.PanicIfError(err)
}

func (repository *UserRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, userId int) (domain.User, error) {
	SQL := `SELECT 
				a.user_id, 
				a.user_name, 
				a.user_email, 
				IFNULL(a.user_token,''), 
				IFNULL(a.user_token_refresh,''), 
				a.user_lang_code, 
				IFNULL(a.user_last_login,''), 
				IFNULL(a.created_by,0), 
				IFNULL(b.user_name,''),
				IFNULL(a.created_at,''), 
				IFNULL(a.updated_by,0), 
				IFNULL(c.user_name,''),
				IFNULL(a.updated_at,'')
			FROM 
				user a 
			LEFT JOIN
				user b ON b.user_id = a.created_by
			LEFT JOIN
				user c ON c.user_id = a.updated_by
			WHERE 
				a.user_id = ?
				AND a.deleted_at IS NULL`
	rows, err := tx.QueryContext(ctx, SQL, userId)
	helper.PanicIfError(err)
	defer rows.Close()

	user := domain.User{}
	if rows.Next() {
		err := rows.Scan(
			&user.UserId,
			&user.UserName,
			&user.UserEmail,
			&user.UserToken,
			&user.UserTokenRefresh,
			&user.UserLangCode,
			&user.UserLastLogin,
			&user.CreatedBy,
			&user.CreatedByName,
			&user.CreatedAt,
			&user.UpdatedBy,
			&user.UpdatedByName,
			&user.UpdatedAt)
		helper.PanicIfError(err)
	}

	return user, nil
}

func (repository *UserRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.User {
	SQL := `SELECT 
				a.user_id, 
				a.user_name, 
				a.user_email, 
				IFNULL(a.user_token,''), 
				IFNULL(a.user_token_refresh,''), 
				a.user_lang_code, 
				IFNULL(a.user_last_login,''), 
				IFNULL(a.created_by,0), 
				IFNULL(b.user_name,''),
				IFNULL(a.created_at,''), 
				IFNULL(a.updated_by,0), 
				IFNULL(c.user_name,''),
				IFNULL(a.updated_at,'')
			FROM 
				user a
			LEFT JOIN
				user b ON b.user_id = a.created_by
			LEFT JOIN
				user c ON c.user_id = a.updated_by
			WHERE
				a.deleted_at IS NULL`
	rows, err := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(err)
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		user := domain.User{}
		err := rows.Scan(
			&user.UserId,
			&user.UserName,
			&user.UserEmail,
			&user.UserToken,
			&user.UserTokenRefresh,
			&user.UserLangCode,
			&user.UserLastLogin,
			&user.CreatedBy,
			&user.CreatedByName,
			&user.CreatedAt,
			&user.UpdatedBy,
			&user.UpdatedByName,
			&user.UpdatedAt)
		helper.PanicIfError(err)
		users = append(users, user)
	}

	return users
}

func (repository *UserRepositoryImpl) FindByEmail(ctx context.Context, tx *sql.Tx, userEmail string) (domain.User, error) {
	SQL := `SELECT 
				user_id, 
				user_name, 
				user_password, 
				user_lang_code 
			FROM 
				user 
			WHERE 
				user_email = ?
				AND deleted_at IS NULL`
	rows, err := tx.QueryContext(ctx, SQL, userEmail)
	helper.PanicIfError(err)
	defer rows.Close()

	user := domain.User{}
	if rows.Next() {
		err := rows.Scan(
			&user.UserId,
			&user.UserName,
			&user.UserPassword,
			&user.UserLangCode)
		helper.PanicIfError(err)
	}

	return user, nil
}

func (repository *UserRepositoryImpl) FindByTokenRefresh(ctx context.Context, tx *sql.Tx, userTokenRefresh string) (domain.User, error) {
	SQL := `SELECT 
				user_id, 
				user_name, 
				user_password, 
				user_lang_code 
			FROM 
				user 
			WHERE 
				user_token_refresh = ?
				AND deleted_at IS NULL`
	rows, err := tx.QueryContext(ctx, SQL, userTokenRefresh)
	helper.PanicIfError(err)
	defer rows.Close()

	user := domain.User{}
	if rows.Next() {
		err := rows.Scan(
			&user.UserId,
			&user.UserName,
			&user.UserPassword,
			&user.UserLangCode)
		helper.PanicIfError(err)
	}

	return user, nil
}

func (repository *UserRepositoryImpl) UpdateToken(ctx context.Context, tx *sql.Tx, user domain.User) domain.User {
	SQL := `UPDATE 
				user 
			SET 
				user_token = ?, 
				user_token_refresh = ?, 
				user_last_login = ? 
			WHERE 
				user_id = ?
				AND deleted_at IS NULL`
	_, err := tx.ExecContext(ctx, SQL,
		user.UserToken,
		user.UserTokenRefresh,
		user.UserLastLogin,
		user.UserId)
	helper.PanicIfError(err)

	return user
}

func (repository *UserRepositoryImpl) Logout(ctx context.Context, tx *sql.Tx, user domain.User) domain.User {
	SQL := `UPDATE 
				user 
			SET 
				user_token = NULL, 
				user_token_refresh = NULL 
			WHERE 
				user_id = ?
				AND deleted_at IS NULL`
	_, err := tx.ExecContext(ctx, SQL, user.UserId)
	helper.PanicIfError(err)

	return user
}
