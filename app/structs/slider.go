package structs

import (
	"time"
)

type Slider struct {
	Title    string    `json:"title" db:"title"`
	Thumb    string    `json:"thumb" db:"thumb"`
	Created  time.Time `json:"created" db:"created"`
	Modified time.Time `json:"modified" db:"modified"`
}
