package structs

import (
	"time"

	"github.com/google/uuid"
)

type Slide struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Title    string    `json:"title" db:"title"`
	Thumb    string    `json:"thumb" db:"thumb"`
	Alt      string    `json:"alt" db:"alt"`
	Created  time.Time `json:"created" db:"created"`
	Modified time.Time `json:"modified" db:"modified"`
}

type SliderResponse struct {
	Code   int     `json:"code"`
	Slider []Slide `json:"slider"`
}
