package web

import (
	"collapp/module/lang/model/domain"
)

// request
type LangCreateRequest struct {
	LangCode  string `validate:"required,min=1,max=5" json:"lang_code"`
	LangName  string `validate:"required,min=1,max=255" json:"lang_name"`
	CreatedBy int    `validate:"required"`
	CreatedAt string `validate:"required"`
}

type LangUpdateRequest struct {
	LangId    int    `validate:"required"`
	LangCode  string `validate:"required,min=1,max=5" json:"lang_code"`
	LangName  string `validate:"required,min=1,max=255" json:"lang_name"`
	UpdatedBy int    `validate:"required"`
	UpdatedAt string `validate:"required"`
}

// rersponse
type LangResponse struct {
	LangId        int    `json:"lang_id"`
	LangCode      string `json:"lang_code"`
	LangName      string `json:"lang_name"`
	CreatedBy     int    `json:"created_by"`
	CreatedByName string `json:"created_by_name"`
	CreatedAt     string `json:"created_at"`
	UpdatedBy     int    `json:"updated_by"`
	UpdatedByName string `json:"updated_by_name"`
	UpdatedAt     string `json:"updated_at"`
}

func ToLangResponse(lang domain.Lang) LangResponse {
	return LangResponse{
		LangId:        lang.LangId,
		LangCode:      lang.LangCode,
		LangName:      lang.LangName,
		CreatedBy:     lang.CreatedBy,
		CreatedByName: lang.CreatedByName,
		CreatedAt:     lang.CreatedAt,
		UpdatedBy:     lang.UpdatedBy,
		UpdatedByName: lang.UpdatedByName,
		UpdatedAt:     lang.UpdatedAt,
	}
}

func ToLangResponses(langs []domain.Lang) []LangResponse {
	var langResponses []LangResponse
	for _, lang := range langs {
		langResponses = append(langResponses, ToLangResponse(lang))
	}
	return langResponses
}
