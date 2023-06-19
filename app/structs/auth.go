package structs

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	AccessCookieName  = "access-cookie"
	RefreshCookieName = "refresh-cookie"
)

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

