package services

import (
	"github.com/Mahaveer86619/FrameSense/pkg/errz"
	"github.com/Mahaveer86619/FrameSense/pkg/middleware"
	"github.com/Mahaveer86619/FrameSense/pkg/models"
	"github.com/Mahaveer86619/FrameSense/pkg/views"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userService *UserService
}

func NewAuthService(userService *UserService) *AuthService {
	return &AuthService{userService: userService}
}

func (s *AuthService) Register(req views.RegisterRequest) (*views.TokenResponse, error) {
	// 1. Check if user exists
	if s.userService.Exists(req.Email, req.Username) {
		return nil, errz.New(errz.Conflict, "Email or Username already exists", nil)
	}

	// 2. Hash Password
	hashedParams, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errz.New(errz.InternalServerError, "Failed to process password", err)
	}

	// 3. Create User
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedParams),
	}

	if err := s.userService.Create(user); err != nil {
		return nil, errz.New(errz.InternalServerError, "Failed to create user", err)
	}

	// 4. Generate Tokens
	accessToken, refreshToken, err := middleware.GenerateTokens(user.ID)
	if err != nil {
		return nil, errz.New(errz.InternalServerError, "Failed to generate tokens", err)
	}

	response := views.ToTokenResponse(accessToken, refreshToken, user)
	return &response, nil
}

func (s *AuthService) Login(req views.LoginRequest) (*views.TokenResponse, error) {
	// 1. Find User
	user, err := s.userService.FindByEmail(req.Email)
	if err != nil {
		return nil, errz.New(errz.Unauthorized, "Invalid credentials", err)
	}

	// 2. Check Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errz.New(errz.Unauthorized, "Invalid credentials", nil)
	}

	// 3. Generate Tokens
	accessToken, refreshToken, err := middleware.GenerateTokens(user.ID)
	if err != nil {
		return nil, errz.New(errz.InternalServerError, "Failed to generate tokens", err)
	}

	response := views.ToTokenResponse(accessToken, refreshToken, user)
	return &response, nil
}

func (s *AuthService) Refresh(req views.RefreshTokenRequest) (map[string]string, error) {
	// 1. Validate Refresh Token
	userID, err := middleware.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// 2. Check if user still exists (optional but recommended)
	if _, err := s.userService.FindByID(userID); err != nil {
		return nil, errz.New(errz.Unauthorized, "User no longer exists", err)
	}

	// 3. Generate New Tokens
	accessToken, newRefreshToken, err := middleware.GenerateTokens(userID)
	if err != nil {
		return nil, errz.New(errz.InternalServerError, "Failed to refresh tokens", err)
	}

	return views.ToTokenOnlyResponse(accessToken, newRefreshToken), nil
}
