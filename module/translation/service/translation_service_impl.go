package service

import (
	"collapp/helper"
	"collapp/module/translation/model"
	"collapp/module/translation/repository"
	"context"
	"database/sql"
)

type TranslationServiceImpl struct {
	TranslationRepository repository.TranslationRepository
	DB                    *sql.DB
}

func NewTranslationService(DB *sql.DB, repo repository.TranslationRepository) TranslationService {
	return &TranslationServiceImpl{
		TranslationRepository: repo,
		DB:                    DB,
	}
}

func (service *TranslationServiceImpl) Create(ctx context.Context, request model.TranslationCreateRequest) model.TranslationResponse {
	tx, err := service.DB.Begin()
	helper.IfError(err)
	defer helper.CommitOrRollback(tx)

	translationData := service.TranslationRepository.Save(ctx, tx, request)
	if translationData.TranslationId > 0 {
		var res = true
		requestText := model.TranslationTextRequest{}

		for _, dt := range request.TranslationText {
			requestText.TranslationTextTranslationId = translationData.TranslationId
			requestText.TranslationTextLangCode = dt.TranslationTextLangCode
			requestText.TranslationTextLangText = dt.TranslationTextLangText
			res = res && service.TranslationRepository.SaveText(ctx, tx, requestText)
		}

		if res {
			translationData, err := service.TranslationRepository.FindById(ctx, tx, translationData.TranslationId)
			helper.IfError(err)

			translationData.TranslationText = service.TranslationRepository.TextFindById(ctx, tx, translationData.TranslationId)
			return model.ToTranslationResponse(translationData)
		} else {
			return model.ToTranslationResponse(translationData)
		}

	} else {
		return model.ToTranslationResponse(translationData)
	}
}

func (service *TranslationServiceImpl) Update(ctx context.Context, request model.TranslationUpdateRequest) model.TranslationResponse {
	tx, err := service.DB.Begin()
	helper.IfError(err)
	defer helper.CommitOrRollback(tx)

	translationData, err := service.TranslationRepository.FindById(ctx, tx, request.TranslationId)
	if err == nil {
		translationData = service.TranslationRepository.Update(ctx, tx, request)

		requestDeleteText := model.TranslationTextDeleteRequest{}
		requestDeleteText.TranslationTextTranslationId = request.TranslationId
		service.TranslationRepository.DeleteText(ctx, tx, requestDeleteText)

		var res = true
		requestText := model.TranslationTextRequest{}

		for _, dt := range request.TranslationText {
			requestText.TranslationTextTranslationId = translationData.TranslationId
			requestText.TranslationTextLangCode = dt.TranslationTextLangCode
			requestText.TranslationTextLangText = dt.TranslationTextLangText
			res = res && service.TranslationRepository.SaveText(ctx, tx, requestText)
		}

		if res {
			translationData, err := service.TranslationRepository.FindById(ctx, tx, translationData.TranslationId)
			helper.IfError(err)

			translationData.TranslationText = service.TranslationRepository.TextFindById(ctx, tx, translationData.TranslationId)
			return model.ToTranslationResponse(translationData)
		} else {
			return model.ToTranslationResponse(translationData)
		}
	}

	return model.ToTranslationResponse(translationData)
}

func (service *TranslationServiceImpl) Delete(ctx context.Context, translationId int) model.TranslationResponse {
	tx, err := service.DB.Begin()
	helper.IfError(err)
	defer helper.CommitOrRollback(tx)

	translationData, err := service.TranslationRepository.FindById(ctx, tx, translationId)
	if err == nil {
		service.TranslationRepository.Delete(ctx, tx, translationData)

		requestDeleteText := model.TranslationTextDeleteRequest{}
		requestDeleteText.TranslationTextTranslationId = translationId
		service.TranslationRepository.DeleteText(ctx, tx, requestDeleteText)
	}

	return model.ToTranslationResponse(translationData)
}

func (service *TranslationServiceImpl) FindById(ctx context.Context, translationId int) model.TranslationResponse {
	tx, err := service.DB.Begin()
	helper.IfError(err)
	defer helper.CommitOrRollback(tx)

	translationData, _ := service.TranslationRepository.FindById(ctx, tx, translationId)
	translationData.TranslationText = service.TranslationRepository.TextFindById(ctx, tx, translationId)

	return model.ToTranslationResponse(translationData)
}

func (service *TranslationServiceImpl) FindAll(ctx context.Context) []model.TranslationResponse {
	tx, err := service.DB.Begin()
	helper.IfError(err)
	defer helper.CommitOrRollback(tx)

	var translationsData = service.TranslationRepository.FindAll(ctx, tx)

	for index, dt := range translationsData {
		translationsData[index].TranslationText = service.TranslationRepository.TextFindById(ctx, tx, dt.TranslationId)
	}

	return model.ToTranslationResponses(translationsData)
}

func (service *TranslationServiceImpl) Translation(ctx context.Context, key string, langCode string) string {
	tx, err := service.DB.Begin()
	helper.IfError(err)
	defer helper.CommitOrRollback(tx)

	translation := service.TranslationRepository.Translation(ctx, tx, key, langCode)

	return translation
}

func (service *TranslationServiceImpl) CheckKeyTranslationExist(ctx context.Context, key string) bool {
	tx, err := service.DB.Begin()
	helper.IfError(err)
	defer helper.CommitOrRollback(tx)

	keyIsExsit := service.TranslationRepository.CheckKeyTranslationExist(ctx, tx, key)

	return keyIsExsit
}
