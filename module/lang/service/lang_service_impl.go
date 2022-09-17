package service

import (
	"collapp/helper"
	"collapp/module/lang/model"
	"collapp/module/lang/repository"
	"context"
	"database/sql"
)

type LangServiceImpl struct {
	LangRepository repository.LangRepository
	DB             *sql.DB
}

func NewLangService(DB *sql.DB) LangService {
	langRepository := repository.NewLangRepository()

	return &LangServiceImpl{
		LangRepository: langRepository,
		DB:             DB,
	}
}

func (service *LangServiceImpl) Create(ctx context.Context, request model.LangCreateRequest) model.LangResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	langData := service.LangRepository.Save(ctx, tx, request)
	if langData.LangId > 0 {
		langData, err := service.LangRepository.FindById(ctx, tx, langData.LangId)
		helper.PanicIfError(err)

		return model.ToLangResponse(langData)
	} else {
		return model.ToLangResponse(langData)
	}
}

func (service *LangServiceImpl) Update(ctx context.Context, request model.LangUpdateRequest) model.LangResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	langData, err := service.LangRepository.FindById(ctx, tx, request.LangId)
	if err == nil {
		langData = service.LangRepository.Update(ctx, tx, request)

		langData, err := service.LangRepository.FindById(ctx, tx, langData.LangId)
		helper.PanicIfError(err)

		return model.ToLangResponse(langData)
	}

	return model.ToLangResponse(langData)
}

func (service *LangServiceImpl) Delete(ctx context.Context, langId int) model.LangResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	langData, err := service.LangRepository.FindById(ctx, tx, langId)
	if err == nil {
		service.LangRepository.Delete(ctx, tx, langData)
	}

	return model.ToLangResponse(langData)
}

func (service *LangServiceImpl) FindById(ctx context.Context, langId int) model.LangResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	langData, _ := service.LangRepository.FindById(ctx, tx, langId)

	return model.ToLangResponse(langData)
}

func (service *LangServiceImpl) FindAll(ctx context.Context) []model.LangResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	langsData := service.LangRepository.FindAll(ctx, tx)

	return model.ToLangResponses(langsData)
}
