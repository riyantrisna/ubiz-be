package repository

import (
	"collapp/module/lang/model"
	"context"
	"database/sql"
)

type LangRepository interface {
	Save(ctx context.Context, tx *sql.Tx, lang model.LangCreateRequest) model.Lang
	Update(ctx context.Context, tx *sql.Tx, lang model.LangUpdateRequest) model.Lang
	Delete(ctx context.Context, tx *sql.Tx, lang model.Lang)
	FindById(ctx context.Context, tx *sql.Tx, langId int) (model.Lang, error)
	FindAll(ctx context.Context, tx *sql.Tx) []model.Lang
}
