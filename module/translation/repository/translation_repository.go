package repository

import (
	"collapp/module/translation/model"
	"context"
	"database/sql"
)

type TranslationRepository interface {
	Save(ctx context.Context, tx *sql.Tx, translation model.TranslationCreateRequest) model.Translation
	SaveText(ctx context.Context, tx *sql.Tx, translation model.TranslationTextRequest) bool
	Update(ctx context.Context, tx *sql.Tx, translation model.TranslationUpdateRequest) model.Translation
	Delete(ctx context.Context, tx *sql.Tx, translation model.Translation)
	DeleteText(ctx context.Context, tx *sql.Tx, translation model.TranslationTextDeleteRequest)
	FindById(ctx context.Context, tx *sql.Tx, translationId int) (model.Translation, error)
	TextFindById(ctx context.Context, tx *sql.Tx, translationId int) []model.TranslationText
	FindAll(ctx context.Context, tx *sql.Tx) []model.Translation
}
