package structs

import (
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Thumb       string    `json:"thumb" db:"thumb"`
	Alt         string    `json:"alt" db:"alt"`
	Description string    `json:"description" db:"description"`
	Created     time.Time `json:"created" db:"created"`
	Modified    time.Time `json:"modified" db:"modified"`
}

type CardsResponse struct {
	Code  int    `json:"code"`
	Tolal int    `json:"total"`
	Data  []Card `json:"data"`
}

// func NewCard(title, thumb, alt, description string) Card {
// 	return Card{
// 		Title:       title,
// 		Thumb:       thumb,
// 		Alt:         alt,
// 		Description: parsDescription(description),
// 	}
// }

// func parsDescription(description string) []string {
// 	str := strings.Trim(description, "[]")

// 	rawSlice := strings.Split(str, ",")

// 	res := []string{}

// 	for _, i := range rawSlice {
// 		res = append(res, strings.Trim(i, "\""))
// 	}

// 	return res
// }
