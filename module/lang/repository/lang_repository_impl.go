package repository

import (
	"collapp/helper"
	"collapp/module/lang/model/domain"
	"context"
	"database/sql"
)

type LangRepositoryImpl struct {
}

func NewLangRepository() LangRepository {
	return &LangRepositoryImpl{}
}

func (repository *LangRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, lang domain.Lang) domain.Lang {

	SQL := `INSERT INTO lang
			(
				lang_code, 
				lang_name, 
				created_by, 
				created_at
			) VALUES (
				?, 
				?, 
				?, 
				?
			)`
	result, err := tx.ExecContext(ctx, SQL,
		lang.LangCode,
		lang.LangName,
		lang.CreatedBy,
		lang.CreatedAt)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	lang.LangId = int(id)
	return lang
}

func (repository *LangRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, lang domain.Lang) domain.Lang {
	SQL := `UPDATE 
				lang 
			SET 
				lang_code = ?, 
				lang_name = ?, 
				updated_by = ?, 
				updated_at = ? 
			WHERE 
				lang_id = ?`
	_, err := tx.ExecContext(ctx, SQL,
		lang.LangCode,
		lang.LangName,
		lang.UpdatedBy,
		lang.UpdatedAt,
		lang.LangId)
	helper.PanicIfError(err)

	return lang
}

func (repository *LangRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, lang domain.Lang) {
	SQL := `DELETE FROM lang WHERE lang_id = ?`
	_, err := tx.ExecContext(ctx, SQL, lang.LangId)
	helper.PanicIfError(err)
}

func (repository *LangRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, langId int) (domain.Lang, error) {
	SQL := `SELECT 
				lang_id, 
				lang_code, 
				lang_name, 
				IFNULL(created_by,0), 
				IFNULL(created_at,''), 
				IFNULL(updated_by,0), 
				IFNULL(updated_at,'') 
			FROM 
				lang 
			WHERE 
				lang_id = ?`
	rows, err := tx.QueryContext(ctx, SQL, langId)
	helper.PanicIfError(err)
	defer rows.Close()

	lang := domain.Lang{}
	if rows.Next() {
		err := rows.Scan(
			&lang.LangId,
			&lang.LangCode,
			&lang.LangName,
			&lang.CreatedBy,
			&lang.CreatedAt,
			&lang.UpdatedBy,
			&lang.UpdatedAt)
		helper.PanicIfError(err)
	}

	return lang, nil
}

func (repository *LangRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.Lang {
	SQL := `SELECT 
				lang_id, 
				lang_code, 
				lang_name, 
				IFNULL(created_by,0), 
				IFNULL(created_at,''), 
				IFNULL(updated_by,0), 
				IFNULL(updated_at,'') 
			FROM 
				lang`
	rows, err := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(err)
	defer rows.Close()

	var langs []domain.Lang
	for rows.Next() {
		lang := domain.Lang{}
		err := rows.Scan(
			&lang.LangId,
			&lang.LangCode,
			&lang.LangName,
			&lang.CreatedBy,
			&lang.CreatedAt,
			&lang.UpdatedBy,
			&lang.UpdatedAt)
		helper.PanicIfError(err)
		langs = append(langs, lang)
	}

	return langs
}
