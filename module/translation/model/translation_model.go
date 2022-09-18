package model

import (
	"database/sql"
)

// model Translation
type Translation struct {
	TranslationId   int
	TranslationKey  string
	TranslationText []TranslationText
	CreatedBy       int
	CreatedByCheck  sql.NullInt32
	CreatedAt       string
	CreatedAtCheck  sql.NullString
	UpdatedBy       int
	UpdatedByCheck  sql.NullInt32
	UpdatedAt       string
	UpdatedAtCheck  sql.NullString
}

type TranslationText struct {
	TranslationTextTranslationId int
	TranslationTextLangCode      string
	TranslationTextLangText      string
}

// request
type TranslationCreateRequest struct {
	TranslationKey  string                   `validate:"required,min=1,max=255" json:"translation_key"`
	TranslationText []TranslationTextRequest `json:"translation_text"`
	CreatedBy       int                      `validate:"required"`
	CreatedAt       string                   `validate:"required"`
}

type TranslationUpdateRequest struct {
	TranslationId   int                      `validate:"required"`
	TranslationKey  string                   `validate:"required,min=1,max=255" json:"translation_key"`
	TranslationText []TranslationTextRequest `json:"translation_text"`
	UpdatedBy       int                      `validate:"required"`
	UpdatedAt       string                   `validate:"required"`
}

type TranslationTextRequest struct {
	TranslationTextTranslationId int    `validate:"required"`
	TranslationTextLangCode      string `validate:"required,min=1,max=255" json:"lang_code"`
	TranslationTextLangText      string `validate:"required" json:"lang_text"`
}

type TranslationTextDeleteRequest struct {
	TranslationTextTranslationId int `validate:"required" json:"lang_translation_id"`
}

// rersponse
type TranslationResponse struct {
	TranslationId   int                       `json:"translation_id"`
	TranslationKey  string                    `json:"translation_code"`
	TranslationText []TranslationTextResponse `json:"translation_text"`
	CreatedBy       int                       `json:"created_by"`
	CreatedAt       string                    `json:"created_at"`
	UpdatedBy       int                       `json:"updated_by"`
	UpdatedAt       string                    `json:"updated_at"`
}

type TranslationTextResponse struct {
	TranslationTextTranslationId int    `json:"lang_translation_id"`
	TranslationTextLangCode      string `json:"lang_code"`
	TranslationTextLangText      string `json:"lang_text"`
}

func ToTranslationResponse(translation Translation) TranslationResponse {
	return TranslationResponse{
		TranslationId:   translation.TranslationId,
		TranslationKey:  translation.TranslationKey,
		TranslationText: ToTranslationTextResponses(translation.TranslationText),
		CreatedBy:       translation.CreatedBy,
		CreatedAt:       translation.CreatedAt,
		UpdatedBy:       translation.UpdatedBy,
		UpdatedAt:       translation.UpdatedAt,
	}
}

func ToTranslationResponses(translations []Translation) []TranslationResponse {
	var translationResponses []TranslationResponse
	for _, translation := range translations {
		translationResponses = append(translationResponses, ToTranslationResponse(translation))
	}
	return translationResponses
}

func ToTranslationTextResponse(text TranslationText) TranslationTextResponse {
	return TranslationTextResponse{
		TranslationTextTranslationId: text.TranslationTextTranslationId,
		TranslationTextLangCode:      text.TranslationTextLangCode,
		TranslationTextLangText:      text.TranslationTextLangText,
	}
}

func ToTranslationTextResponses(texts []TranslationText) []TranslationTextResponse {
	var textResponses []TranslationTextResponse
	for _, text := range texts {
		textResponses = append(textResponses, ToTranslationTextResponse(text))
	}
	return textResponses
}
