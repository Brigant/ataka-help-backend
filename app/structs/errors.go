package structs

import "errors"

var (
	ErrForeignViolation     = errors.New("wrong foreign key")
	ErrUniqueRestriction    = errors.New("violation of the database unique restriction")
	ErrNoRowAffected        = errors.New("no one row was affected")
	ErrNotFound             = errors.New("nothing was found")
	ErrCheckCaptcha         = errors.New("captha doesn't confirm the checking")
	ErrMailSending          = errors.New("error happaned while sendin email")
	ErrWrongBody            = errors.New("wrong body")
	ErrInvalidSigningMethod = errors.New("invalid signing metod")
	ErrWrongTokenClaimType  = errors.New("token claims are not of type *tokenClaims")
	ErrTimeout              = errors.New("timeout")
	ErrEmptyField           = errors.New("field shoud be not empty")
	ErrNotMatch             = errors.New("the pasword doesn't match")
	ErrNoSession            = errors.New("there is no such session")
)
