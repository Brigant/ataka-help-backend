package structs

import (
	"time"

	"github.com/google/uuid"
)

type Partner struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Alt      string    `json:"alt" db:"alt"`
	Thumb    string    `json:"thumb" db:"thumb"`
	Created  time.Time `json:"created" db:"created"`
	Modified time.Time `json:"modified" db:"modified"`
}

type PartnerResponse struct {
	Code     int       `json:"code"`
	Total    int       `json:"total"`
	Partners []Partner `json:"partners"`
}

type PartnerQueryParameters struct {
	Offset int `json:"offset" query:"offset"`
	Limit  int `json:"limit" query:"limit"`
}
