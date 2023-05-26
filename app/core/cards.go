package core

import (
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID           uuid.UUID     `json:"id"`
	Title        string        `json:"title" form:"title" validate:"required,min=3,max=32"`
	Image        string        `json:"image"`
	Alt          string        `json:"alt" form:"alt"`
	Descriptions []Description `json:"description" form:""`
	Created      time.Time     `json:"created"`
	Modified     time.Time     `json:"modified"`
}

type Description struct {
	Item string
}
