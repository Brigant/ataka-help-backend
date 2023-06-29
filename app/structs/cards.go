package structs

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID          uuid.UUID       `json:"id" db:"id" params:"id"`
	Title       string          `json:"title" db:"title"`
	Thumb       string          `json:"thumb" db:"thumb"`
	Alt         string          `json:"alt" db:"alt"`
	Description json.RawMessage `json:"description" db:"description"`
	Created     time.Time       `json:"created" db:"created"`
	Modified    time.Time       `json:"modified" db:"modified"`
}

type CardsResponse struct {
	Code  int    `json:"code"`
	Total int    `json:"total"`
	Cards []Card `json:"cards"`
}

type CardQueryParameters struct {
	Page  int `json:"page" query:"page"`
	Limit int `json:"limit" query:"limit"`
}
