package views

import (
	"time"

	"github.com/Mahaveer86619/FrameSense/pkg/models"
	"github.com/Mahaveer86619/FrameSense/pkg/utils"
)

// --- Requests ---

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// --- Responses ---

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type TokenResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

// --- Mappers ---

func ToUserResponse(u *models.User) UserResponse {
	return UserResponse{
		ID:        utils.MaskID(u.ID),
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

func ToTokenResponse(accessToken, refreshToken string, user *models.User) TokenResponse {
	return TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         ToUserResponse(user),
	}
}

func ToTokenOnlyResponse(accessToken, refreshToken string) map[string]string {
	return map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
}
