package structs

import "errors"

var (
	ErrForeignViolation  = errors.New("wrong foreign key")
	ErrUniqueRestriction = errors.New("violation of the database unique restriction")
	ErrNoRowAffected     = errors.New("no one row was affected")
	ErrNotFound          = errors.New("nothing was found")
	ErrCheckCaptcha      = errors.New("captha doesn't confirm the checking")
	ErrMailSending       = errors.New("error happaned while sendin email")
)
