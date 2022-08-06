package service

import (
	"collapp/helper"
	"collapp/module/lang/model/domain"
	"collapp/module/lang/model/web"
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

func (service *LangServiceImpl) Create(ctx context.Context, request web.LangCreateRequest) web.LangResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	langData := domain.Lang{
		LangCode:  request.LangCode,
		LangName:  request.LangName,
		CreatedBy: request.CreatedBy,
		CreatedAt: request.CreatedAt,
	}

	langData = service.LangRepository.Save(ctx, tx, langData)

	return web.ToLangResponse(langData)
}

func (service *LangServiceImpl) Update(ctx context.Context, request web.LangUpdateRequest) web.LangResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	langData, err := service.LangRepository.FindById(ctx, tx, request.LangId)
	if err == nil {
		langData.LangCode = request.LangCode
		langData.LangName = request.LangName
		langData.UpdatedBy = request.UpdatedBy
		langData.UpdatedAt = request.UpdatedAt

		langData = service.LangRepository.Update(ctx, tx, langData)
	}

	return web.ToLangResponse(langData)
}

func (service *LangServiceImpl) Delete(ctx context.Context, langId int) web.LangResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	langData, err := service.LangRepository.FindById(ctx, tx, langId)
	if err == nil {
		service.LangRepository.Delete(ctx, tx, langData)
	}

	return web.ToLangResponse(langData)
}

func (service *LangServiceImpl) FindById(ctx context.Context, langId int) web.LangResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	langData, _ := service.LangRepository.FindById(ctx, tx, langId)

	return web.ToLangResponse(langData)
}

func (service *LangServiceImpl) FindAll(ctx context.Context) []web.LangResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	langsData := service.LangRepository.FindAll(ctx, tx)

	return web.ToLangResponses(langsData)
}
