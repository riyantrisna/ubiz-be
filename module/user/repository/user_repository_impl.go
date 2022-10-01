package repository

import (
	"collapp/helper"
	"collapp/module/user/model"
	"context"
	"database/sql"
)

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, user model.UserCreateRequest) model.User {

	SQL := `INSERT INTO user
			(
				user_name, 
				user_email, 
				user_password, 
				user_lang_code, 
				user_photo,
				created_by, 
				created_at
			) VALUES (
				?, 
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
		user.UserPhotoName,
		user.CreatedBy,
		user.CreatedAt)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	res := model.User{}
	res.UserId = int(id)
	return res
}

func (repository *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user model.UserUpdateRequest) model.User {
	SQL := `UPDATE 
				user 
			SET 
				user_name = ?, 
				user_email = ?, 
				user_lang_code = ?, 
				user_photo = ?,
				updated_by = ?, 
				updated_at = ? 
			WHERE 
				user_id = ?`
	_, err := tx.ExecContext(ctx, SQL,
		user.UserName,
		user.UserEmail,
		user.UserLangCode,
		user.UserPhotoName,
		user.UpdatedBy,
		user.UpdatedAt,
		user.UserId)
	helper.PanicIfError(err)

	res := model.User{}
	res.UserId = user.UserId
	return res
}

func (repository *UserRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, user model.User) {
	SQL := `DELETE FROM user WHERE user_id = ?`
	_, err := tx.ExecContext(ctx, SQL, user.UserId)
	helper.PanicIfError(err)
}

func (repository *UserRepositoryImpl) SoftDelete(ctx context.Context, tx *sql.Tx, user model.User) {
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

func (repository *UserRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, userId int) (model.User, error) {
	SQL := `SELECT 
				a.user_id, 
				a.user_name, 
				a.user_email, 
				a.user_token, 
				a.user_token_refresh, 
				a.user_lang_code, 
				a.user_last_login, 
				a.user_photo,
				a.created_by, 
				b.user_name,
				a.created_at, 
				a.updated_by,
				c.user_name,
				a.updated_at
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

	user := model.User{}
	if rows.Next() {
		err := rows.Scan(
			&user.UserId,
			&user.UserName,
			&user.UserEmail,
			&user.UserTokenCheck,
			&user.UserTokenRefreshCheck,
			&user.UserLangCode,
			&user.UserLastLoginCheck,
			&user.UserPhotoCheck,
			&user.CreatedByCheck,
			&user.CreatedByNameCheck,
			&user.CreatedAtCheck,
			&user.UpdatedByCheck,
			&user.UpdatedByNameCheck,
			&user.UpdatedAtCheck)
		helper.PanicIfError(err)
	}

	if user.UserTokenCheck.Valid {
		user.UserToken = user.UserTokenCheck.String
	}
	if user.UserTokenRefreshCheck.Valid {
		user.UserTokenRefresh = user.UserTokenRefreshCheck.String
	}
	if user.UserLastLoginCheck.Valid {
		user.UserLastLogin = user.UserLastLoginCheck.String
	}
	if user.UserPhotoCheck.Valid {
		user.UserPhoto = user.UserPhotoCheck.String
	}
	if user.CreatedByCheck.Valid {
		user.CreatedBy = int(user.CreatedByCheck.Int32)
	}
	if user.CreatedByNameCheck.Valid {
		user.CreatedByName = user.CreatedByNameCheck.String
	}
	if user.CreatedAtCheck.Valid {
		user.CreatedAt = user.CreatedAtCheck.String
	}
	if user.UpdatedByCheck.Valid {
		user.UpdatedBy = int(user.UpdatedByCheck.Int32)
	}
	if user.UpdatedByNameCheck.Valid {
		user.UpdatedByName = user.UpdatedByNameCheck.String
	}
	if user.UpdatedAtCheck.Valid {
		user.UpdatedAt = user.UpdatedAtCheck.String
	}

	return user, nil
}

func (repository *UserRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []model.User {
	SQL := `SELECT 
				a.user_id, 
				a.user_name, 
				a.user_email, 
				a.user_token, 
				a.user_token_refresh, 
				a.user_lang_code, 
				a.user_last_login, 
				a.user_photo, 
				a.created_by,
				b.user_name,
				a.created_at, 
				a.updated_by,
				c.user_name,
				a.updated_at
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

	var users []model.User
	for rows.Next() {
		user := model.User{}
		err := rows.Scan(
			&user.UserId,
			&user.UserName,
			&user.UserEmail,
			&user.UserTokenCheck,
			&user.UserTokenRefreshCheck,
			&user.UserLangCode,
			&user.UserLastLoginCheck,
			&user.UserPhotoCheck,
			&user.CreatedByCheck,
			&user.CreatedByNameCheck,
			&user.CreatedAtCheck,
			&user.UpdatedByCheck,
			&user.UpdatedByNameCheck,
			&user.UpdatedAtCheck)
		helper.PanicIfError(err)

		if user.UserTokenCheck.Valid {
			user.UserToken = user.UserTokenCheck.String
		}
		if user.UserTokenRefreshCheck.Valid {
			user.UserTokenRefresh = user.UserTokenRefreshCheck.String
		}
		if user.UserLastLoginCheck.Valid {
			user.UserLastLogin = user.UserLastLoginCheck.String
		}
		if user.UserPhotoCheck.Valid {
			user.UserPhoto = user.UserPhotoCheck.String
		}
		if user.CreatedByCheck.Valid {
			user.CreatedBy = int(user.CreatedByCheck.Int32)
		}
		if user.CreatedByNameCheck.Valid {
			user.CreatedByName = user.CreatedByNameCheck.String
		}
		if user.CreatedAtCheck.Valid {
			user.CreatedAt = user.CreatedAtCheck.String
		}
		if user.UpdatedByCheck.Valid {
			user.UpdatedBy = int(user.UpdatedByCheck.Int32)
		}
		if user.UpdatedByNameCheck.Valid {
			user.UpdatedByName = user.UpdatedByNameCheck.String
		}
		if user.UpdatedAtCheck.Valid {
			user.UpdatedAt = user.UpdatedAtCheck.String
		}

		users = append(users, user)
	}

	return users
}

func (repository *UserRepositoryImpl) FindByEmail(ctx context.Context, tx *sql.Tx, userEmail string) (model.User, error) {
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

	user := model.User{}
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

func (repository *UserRepositoryImpl) FindByTokenRefresh(ctx context.Context, tx *sql.Tx, userTokenRefresh string) (model.User, error) {
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

	user := model.User{}
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

func (repository *UserRepositoryImpl) UpdateToken(ctx context.Context, tx *sql.Tx, user model.User) model.User {
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

func (repository *UserRepositoryImpl) Logout(ctx context.Context, tx *sql.Tx, user model.User) model.User {
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
