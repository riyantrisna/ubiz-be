package service

import (
	"collapp/module/lang/model/web"
	"context"
)

type LangService interface {
	Create(ctx context.Context, request web.LangCreateRequest) web.LangResponse
	Update(ctx context.Context, request web.LangUpdateRequest) web.LangResponse
	Delete(ctx context.Context, langId int) web.LangResponse
	FindById(ctx context.Context, langId int) web.LangResponse
	FindAll(ctx context.Context) []web.LangResponse
}
