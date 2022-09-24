package service

import (
	"collapp/module/translation/model"
	"context"
)

type TranslationService interface {
	Create(ctx context.Context, request model.TranslationCreateRequest) model.TranslationResponse
	Update(ctx context.Context, request model.TranslationUpdateRequest) model.TranslationResponse
	Delete(ctx context.Context, translationId int) model.TranslationResponse
	FindById(ctx context.Context, translationId int) model.TranslationResponse
	FindAll(ctx context.Context) []model.TranslationResponse
	Translation(ctx context.Context, key string, langCode string) string
}
