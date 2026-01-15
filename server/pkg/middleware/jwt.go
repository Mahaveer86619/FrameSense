package middleware

import (
	"strings"
	"time"

	"github.com/Mahaveer86619/FrameSense/pkg/config"
	"github.com/Mahaveer86619/FrameSense/pkg/errz"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var jwtKey = []byte(config.AppConfig.JWT_SECRET)

// Use a string constant for Echo context keys
const userIDKey = "user_id"

type Claims struct {
	UserID    uint   `json:"user_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func GenerateTokens(userID uint) (string, string, error) {
	accessExpiration := time.Now().Add(15 * time.Minute)
	accessClaims := &Claims{
		UserID:    userID,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiration),
			Issuer:    "bookture-server",
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshExpiration := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := &Claims{
		UserID:    userID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiration),
			Issuer:    "bookture-server",
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func ValidateRefreshToken(tokenString string) (uint, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return 0, errz.New(errz.Unauthorized, "Invalid refresh token", err)
	}

	if claims.TokenType != "refresh" {
		return 0, errz.New(errz.Unauthorized, "Invalid token type", nil)
	}

	return claims.UserID, nil
}

// Middleware validates the JWT and sets the user_id in the Echo context
func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			errz.HandleErrors(c, errz.New(errz.Unauthorized, "Authorization header required", nil))
			return nil
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			errz.HandleErrors(c, errz.New(errz.Unauthorized, "Invalid or expired token", err))
			return nil
		}

		if claims.TokenType != "access" {
			errz.HandleErrors(c, errz.New(errz.Unauthorized, "Invalid token type for authentication", nil))
			return nil
		}

		// Set UserID in Echo context
		c.Set(userIDKey, claims.UserID)
		return next(c)
	}
}

// OptionalAuth attempts to validate the JWT but continues even if invalid/missing
func OptionalAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader != "" {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
			if err == nil && token.Valid && claims.TokenType == "access" {
				c.Set(userIDKey, claims.UserID)
				return next(c)
			}
		}
		// Continue without user ID if invalid/missing
		return next(c)
	}
}

// GetUserID retrieves the UserID from the Echo context
func GetUserID(c echo.Context) uint {
	val := c.Get(userIDKey)
	if val == nil {
		return 0
	}

	if id, ok := val.(uint); ok {
		return id
	}
	return 0
}
