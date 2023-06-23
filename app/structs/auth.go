package structs

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	AccessCookieName  = "access-cookie"
	RefreshCookieName = "refresh-cookie"
	passMask          = `^[[:graph:]]{8,256}$`
)

var passRegex = regexp.MustCompile(passMask)

type TokenPair struct {
	AccessToken   string `json:"access_token"`
	AccessExpire  time.Time
	RefreshToken  string `json:"refresh_token"`
	RefresgExpire time.Time
}

type IdentityData struct {
	Login    string `json:"login"`
	Password string `jsdon:"password"`
}

type Claims struct {
	jwt.StandardClaims
	UserID string
}

func SHA256(password, salt string) string {
	sum := sha256.Sum256([]byte(password + salt))

	return fmt.Sprintf("%x", sum)
}

type PasswordsContainer struct {
	CurrentPassword      string `json:"currentPassword"`
	NewPassword          string `json:"newPassword"`
	NewPasswordConfirmed string `json:"newPasswordConfirmed"`
}

func (p PasswordsContainer) Validate() error {
	if p.CurrentPassword == "" {
		return fmt.Errorf("old password: %w", ErrEmptyField)
	}

	if p.NewPassword == "" {
		return fmt.Errorf("new password: %w", ErrEmptyField)
	}

	if valid := passRegex.MatchString(p.NewPassword); !valid {
		return fmt.Errorf("password does not match with regex: `%s`", passMask)
	}

	if p.NewPassword != p.NewPasswordConfirmed {
		return ErrNotMatch
	}

	return nil
}
