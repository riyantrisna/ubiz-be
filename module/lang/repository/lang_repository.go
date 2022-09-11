package repository

import (
	"collapp/module/lang/model"
	"context"
	"database/sql"
)

type LangRepository interface {
	Save(ctx context.Context, tx *sql.Tx, lang model.LangCreateRequest) model.LangResponse
	Update(ctx context.Context, tx *sql.Tx, lang model.LangUpdateRequest) model.LangResponse
	Delete(ctx context.Context, tx *sql.Tx, lang model.LangResponse)
	FindById(ctx context.Context, tx *sql.Tx, langId int) (model.LangResponse, error)
	FindAll(ctx context.Context, tx *sql.Tx) []model.LangResponse
}
