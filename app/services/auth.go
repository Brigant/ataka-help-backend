package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/baza-trainee/ataka-help-backend/app/config"
	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/golang-jwt/jwt"
)

type UserRepo interface {
	FindEmailWithPasword(context.Context, structs.IdentityData) (string, error)
}

type AuthService struct {
	Repo UserRepo
}

var inMemory = make(map[string]string)

func (s AuthService) GetTokenPair(ctx context.Context, identity structs.IdentityData, cfg config.AuthConfig) (string, time.Time, error) {
	identity.Password = structs.SHA256(identity.Password, cfg.Salt)

	userID, err := s.Repo.FindEmailWithPasword(ctx, identity)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error in FindEmailWithPasword: %w", err)
	}

	accessExpire := time.Now().Add(cfg.AccessTokenTTL)
	refreshExpire := time.Now().Add(cfg.RefreshTokenTTL)

	accessToken, err := generateJWT(accessExpire, cfg.SigningKey, userID)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error in generateJWT: %w", err)
	}

	refreshToken, err := generateJWT(refreshExpire, cfg.SigningKey, userID)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error in generateJWT: %w", err)
	}

	inMemory[refreshToken] = userID

	tokenPair := structs.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	byteResult, err := json.Marshal(tokenPair)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error in Marshal: %w", err)
	}
	fmt.Println(inMemory)
	return string(byteResult), accessExpire, nil
}

func generateJWT(expire time.Time, signingKey, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, structs.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: userID,
	})

	tokenValue, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", fmt.Errorf("cannot get SignetString token: %w", err)
	}

	return tokenValue, nil
}

func (s AuthService) CleanSession(userID string) {
	for key, value := range inMemory {
		if value == userID {
			delete(inMemory, key)
		}
	}
}

func (s AuthService) Refresh(refreshString string, cfg config.AuthConfig) (string, time.Time, error) {
	userID := inMemory["refreshString"]
	

	delete(inMemory, refreshString)

	accessExpire := time.Now().Add(cfg.AccessTokenTTL)
	refreshExpire := time.Now().Add(cfg.RefreshTokenTTL)

	accessToken, err := generateJWT(accessExpire, cfg.SigningKey, userID)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error in generateJWT: %w", err)
	}

	refreshToken, err := generateJWT(refreshExpire, cfg.SigningKey, userID)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error in generateJWT: %w", err)
	}

	inMemory[refreshToken] = userID

	tokenPair := structs.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	byteResult, err := json.Marshal(tokenPair)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error in Marshal: %w", err)
	}
	fmt.Println(inMemory)
	return string(byteResult), accessExpire, nil
}
