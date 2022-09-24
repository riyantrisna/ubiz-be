package repository

import (
	"collapp/helper"
	"collapp/module/translation/model"
	"context"
	"database/sql"
)

type TranslationRepositoryImpl struct {
}

func NewTranslationRepository() TranslationRepository {
	return &TranslationRepositoryImpl{}
}

func (repository *TranslationRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, translation model.TranslationCreateRequest) model.Translation {

	SQL := `INSERT INTO lang_key
			(
				langkey_key,
				created_by, 
				created_at
			) VALUES (
				?, 
				?, 
				?
			)`
	result, err := tx.ExecContext(ctx, SQL,
		translation.TranslationKey,
		translation.CreatedBy,
		translation.CreatedAt)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	res := model.Translation{}
	res.TranslationId = int(id)
	return res
}

func (repository *TranslationRepositoryImpl) SaveText(ctx context.Context, tx *sql.Tx, translationText model.TranslationTextRequest) bool {

	SQL := `INSERT INTO lang_key_text
			(
				langkeytext_langkey_id, 
				langkeytext_lang_code, 
				langkeytext_lang_text
			) VALUES (
				?, 
				?, 
				?
			)`
	result, err := tx.ExecContext(ctx, SQL,
		translationText.TranslationTextTranslationId,
		translationText.TranslationTextLangCode,
		translationText.TranslationTextLangText)
	helper.PanicIfError(err)

	total, err := result.RowsAffected()
	helper.PanicIfError(err)

	if total > 0 {
		return true
	} else {
		return false
	}
}

func (repository *TranslationRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, translation model.TranslationUpdateRequest) model.Translation {
	SQL := `UPDATE 
				lang_key 
			SET 
				langkey_key = ?,
				updated_by = ?, 
				updated_at = ? 
			WHERE 
				langkey_id = ?`
	_, err := tx.ExecContext(ctx, SQL,
		translation.TranslationKey,
		translation.UpdatedBy,
		translation.UpdatedAt,
		translation.TranslationId)
	helper.PanicIfError(err)

	res := model.Translation{}
	res.TranslationId = translation.TranslationId
	return res
}

func (repository *TranslationRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, translation model.Translation) {
	SQL := `DELETE FROM lang_key WHERE langkey_id = ?`
	_, err := tx.ExecContext(ctx, SQL, translation.TranslationId)
	helper.PanicIfError(err)
}

func (repository *TranslationRepositoryImpl) DeleteText(ctx context.Context, tx *sql.Tx, translation model.TranslationTextDeleteRequest) {
	SQL := `DELETE FROM lang_key_text WHERE langkeytext_langkey_id = ?`
	_, err := tx.ExecContext(ctx, SQL, translation.TranslationTextTranslationId)
	helper.PanicIfError(err)
}

func (repository *TranslationRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, translationId int) (model.Translation, error) {
	SQL := `SELECT 
				langkey_id, 
				langkey_key,
				created_by, 
				created_at, 
				updated_by, 
				updated_at 
			FROM 
				lang_key 
			WHERE 
				langkey_id = ?`
	rows, err := tx.QueryContext(ctx, SQL, translationId)
	helper.PanicIfError(err)
	defer rows.Close()

	translation := model.Translation{}
	if rows.Next() {
		err := rows.Scan(
			&translation.TranslationId,
			&translation.TranslationKey,
			&translation.CreatedByCheck,
			&translation.CreatedAtCheck,
			&translation.UpdatedByCheck,
			&translation.UpdatedAtCheck)
		helper.PanicIfError(err)
	}

	if translation.CreatedByCheck.Valid {
		translation.CreatedBy = int(translation.CreatedByCheck.Int32)
	}
	if translation.CreatedAtCheck.Valid {
		translation.CreatedAt = translation.CreatedAtCheck.String
	}
	if translation.UpdatedByCheck.Valid {
		translation.UpdatedBy = int(translation.UpdatedByCheck.Int32)
	}
	if translation.UpdatedAtCheck.Valid {
		translation.UpdatedAt = translation.UpdatedAtCheck.String
	}

	return translation, nil
}

func (repository *TranslationRepositoryImpl) TextFindById(ctx context.Context, tx *sql.Tx, translationId int) []model.TranslationText {
	SQL := `SELECT 
				langkeytext_langkey_id, 
				langkeytext_lang_code, 
				langkeytext_lang_text
			FROM 
				lang_key_text
			WHERE
				langkeytext_langkey_id = ?`
	rows, err := tx.QueryContext(ctx, SQL, translationId)
	helper.PanicIfError(err)
	defer rows.Close()

	var translations []model.TranslationText
	for rows.Next() {
		translation := model.TranslationText{}
		err := rows.Scan(
			&translation.TranslationTextTranslationId,
			&translation.TranslationTextLangCode,
			&translation.TranslationTextLangText)
		helper.PanicIfError(err)

		translations = append(translations, translation)
	}

	return translations
}

func (repository *TranslationRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []model.Translation {
	SQL := `SELECT 
				langkey_id, 
				langkey_key,
				created_by,
				created_at, 
				updated_by, 
				updated_at 
			FROM 
				lang_key`
	rows, err := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(err)
	defer rows.Close()

	var translations []model.Translation
	for rows.Next() {
		translation := model.Translation{}
		err := rows.Scan(
			&translation.TranslationId,
			&translation.TranslationKey,
			&translation.CreatedByCheck,
			&translation.CreatedAtCheck,
			&translation.UpdatedByCheck,
			&translation.UpdatedAtCheck)
		helper.PanicIfError(err)

		if translation.CreatedByCheck.Valid {
			translation.CreatedBy = int(translation.CreatedByCheck.Int32)
		}
		if translation.CreatedAtCheck.Valid {
			translation.CreatedAt = translation.CreatedAtCheck.String
		}
		if translation.UpdatedByCheck.Valid {
			translation.UpdatedBy = int(translation.UpdatedByCheck.Int32)
		}
		if translation.UpdatedAtCheck.Valid {
			translation.UpdatedAt = translation.UpdatedAtCheck.String
		}

		translations = append(translations, translation)
	}

	return translations
}

func (repository *TranslationRepositoryImpl) Translation(ctx context.Context, tx *sql.Tx, key string, langCode string) string {
	SQL := `SELECT
				a.langkeytext_lang_text
			FROM 
				lang_key_text a
			LEFT JOIN lang_key b ON b.langkey_id = a.langkeytext_langkey_id AND a.langkeytext_lang_code = ?
			WHERE
				b.langkey_key = ?`
	rows, err := tx.QueryContext(ctx, SQL, langCode, key)
	helper.PanicIfError(err)
	defer rows.Close()

	if rows.Next() {
		var translation string
		err := rows.Scan(&translation)
		helper.PanicIfError(err)
		if err == nil {
			return translation
		} else {
			return "[" + key + "]"
		}
	} else {
		return "[" + key + "]"
	}
}

func (repository *TranslationRepositoryImpl) CheckKeyTranslationExist(ctx context.Context, tx *sql.Tx, key string) bool {
	SQL := `SELECT 
				langkey_id
			FROM 
				lang_key
			WHERE
				langkey_key = ?`
	rows, err := tx.QueryContext(ctx, SQL, key)
	helper.PanicIfError(err)
	defer rows.Close()

	if rows.Next() {
		return true
	} else {
		return false
	}
}
