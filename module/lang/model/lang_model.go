package model

import "database/sql"

// model Language
type Lang struct {
	LangId         int
	LangCode       string
	LangName       string
	CreatedBy      int
	CreatedByCheck sql.NullInt32
	CreatedAt      string
	CreatedAtCheck sql.NullString
	UpdatedBy      int
	UpdatedByCheck sql.NullInt32
	UpdatedAt      string
	UpdatedAtCheck sql.NullString
}

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
	LangId    int    `json:"lang_id"`
	LangCode  string `json:"lang_code"`
	LangName  string `json:"lang_name"`
	CreatedBy int    `json:"created_by"`
	CreatedAt string `json:"created_at"`
	UpdatedBy int    `json:"updated_by"`
	UpdatedAt string `json:"updated_at"`
}

func ToLangResponse(lang Lang) LangResponse {
	return LangResponse{
		LangId:    lang.LangId,
		LangCode:  lang.LangCode,
		LangName:  lang.LangName,
		CreatedBy: lang.CreatedBy,
		CreatedAt: lang.CreatedAt,
		UpdatedBy: lang.UpdatedBy,
		UpdatedAt: lang.UpdatedAt,
	}
}

func ToLangResponses(langs []Lang) []LangResponse {
	var langResponses []LangResponse
	for _, lang := range langs {
		langResponses = append(langResponses, ToLangResponse(lang))
	}
	return langResponses
}
