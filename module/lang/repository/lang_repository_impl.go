package repository

import (
	"collapp/helper"
	"collapp/module/lang/model"
	"context"
	"database/sql"
)

type LangRepositoryImpl struct {
}

func NewLangRepository() LangRepository {
	return &LangRepositoryImpl{}
}

func (repository *LangRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, lang model.LangCreateRequest) model.Lang {

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

	res := model.Lang{}
	res.LangId = int(id)
	return res
}

func (repository *LangRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, lang model.LangUpdateRequest) model.Lang {
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

	res := model.Lang{}
	res.LangId = lang.LangId
	return res
}

func (repository *LangRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, lang model.Lang) {
	SQL := `DELETE FROM lang WHERE lang_id = ?`
	_, err := tx.ExecContext(ctx, SQL, lang.LangId)
	helper.PanicIfError(err)
}

func (repository *LangRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, langId int) (model.Lang, error) {
	SQL := `SELECT 
				lang_id, 
				lang_code, 
				lang_name, 
				created_by, 
				created_at, 
				updated_by, 
				updated_at 
			FROM 
				lang 
			WHERE 
				lang_id = ?`
	rows, err := tx.QueryContext(ctx, SQL, langId)
	helper.PanicIfError(err)
	defer rows.Close()

	lang := model.Lang{}
	if rows.Next() {
		err := rows.Scan(
			&lang.LangId,
			&lang.LangCode,
			&lang.LangName,
			&lang.CreatedByCheck,
			&lang.CreatedAtCheck,
			&lang.UpdatedByCheck,
			&lang.UpdatedAtCheck)
		helper.PanicIfError(err)
	}

	if lang.CreatedByCheck.Valid {
		lang.CreatedBy = int(lang.CreatedByCheck.Int32)
	}
	if lang.CreatedAtCheck.Valid {
		lang.CreatedAt = lang.CreatedAtCheck.String
	}
	if lang.UpdatedByCheck.Valid {
		lang.UpdatedBy = int(lang.UpdatedByCheck.Int32)
	}
	if lang.UpdatedAtCheck.Valid {
		lang.UpdatedAt = lang.UpdatedAtCheck.String
	}

	return lang, nil
}

func (repository *LangRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []model.Lang {
	SQL := `SELECT 
				lang_id, 
				lang_code, 
				lang_name, 
				created_by,
				created_at, 
				updated_by, 
				updated_at 
			FROM 
				lang`
	rows, err := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(err)
	defer rows.Close()

	var langs []model.Lang
	for rows.Next() {
		lang := model.Lang{}
		err := rows.Scan(
			&lang.LangId,
			&lang.LangCode,
			&lang.LangName,
			&lang.CreatedByCheck,
			&lang.CreatedAtCheck,
			&lang.UpdatedByCheck,
			&lang.UpdatedAtCheck)
		helper.PanicIfError(err)

		if lang.CreatedByCheck.Valid {
			lang.CreatedBy = int(lang.CreatedByCheck.Int32)
		}
		if lang.CreatedAtCheck.Valid {
			lang.CreatedAt = lang.CreatedAtCheck.String
		}
		if lang.UpdatedByCheck.Valid {
			lang.UpdatedBy = int(lang.UpdatedByCheck.Int32)
		}
		if lang.UpdatedAtCheck.Valid {
			lang.UpdatedAt = lang.UpdatedAtCheck.String
		}

		langs = append(langs, lang)
	}

	return langs
}
