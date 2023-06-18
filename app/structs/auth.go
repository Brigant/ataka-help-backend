package structs

import (
	"crypto/sha256"
	"fmt"

	"github.com/golang-jwt/jwt"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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
