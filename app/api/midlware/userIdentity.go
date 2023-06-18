package midlware

import (
	"encoding/json"
	"fmt"

	"github.com/baza-trainee/ataka-help-backend/app/config"
	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func NewUserIdentity(cfg config.AuthConfig) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tokenString := ctx.Cookies("TokenPair")

		if tokenString == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "empty cooki")
		}

		tokenPair := structs.TokenPair{}

		if err := json.Unmarshal([]byte(tokenString), &tokenPair); err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		userID, err := parseToken(tokenPair.AccessToken, cfg.SigningKey)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		ctx.Locals("userID", userID)
		ctx.Locals("refreshString", tokenPair.RefreshToken)

		return ctx.Next()
	}
}

func parseToken(tokenString, signingKey string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &structs.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, structs.ErrInvalidSigningMethod
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("accessToken throws an error during parsing: %w", err)
	}

	claims, ok := token.Claims.(*structs.Claims)
	if !ok {
		return "", structs.ErrWrongTokenClaimType
	}

	return claims.UserID, nil
}
