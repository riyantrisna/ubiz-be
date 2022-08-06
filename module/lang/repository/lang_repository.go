package repository

import (
	"collapp/module/lang/model/domain"
	"context"
	"database/sql"
)

type LangRepository interface {
	Save(ctx context.Context, tx *sql.Tx, lang domain.Lang) domain.Lang
	Update(ctx context.Context, tx *sql.Tx, lang domain.Lang) domain.Lang
	Delete(ctx context.Context, tx *sql.Tx, lang domain.Lang)
	FindById(ctx context.Context, tx *sql.Tx, langId int) (domain.Lang, error)
	FindAll(ctx context.Context, tx *sql.Tx) []domain.Lang
}
