package structs

type Contact struct {
	Phone1 string `json:"phone1" db:"phone1"`
	Phone2 string `json:"phone2" db:"phone2"`
	Email  string `json:"email" db:"email"`
}
