package structs

import "errors"

var (
	ErrForeignViolation  = errors.New("wrong foreign key")
	ErrUniqueRestriction = errors.New("violation of the database unique restriction")
	ErrDatabaseInserting = errors.New("ntohing is inserted to database")
	ErrNotFound          = errors.New("nothing was found")
	ErrCheckCaptcha      = errors.New("captha doesn't confirm the checking")
)
