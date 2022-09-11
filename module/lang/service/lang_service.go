package service

import (
	"collapp/module/lang/model"
	"context"
)

type LangService interface {
	Create(ctx context.Context, request model.LangCreateRequest) model.LangResponse
	Update(ctx context.Context, request model.LangUpdateRequest) model.LangResponse
	Delete(ctx context.Context, langId int) model.LangResponse
	FindById(ctx context.Context, langId int) model.LangResponse
	FindAll(ctx context.Context) []model.LangResponse
}
