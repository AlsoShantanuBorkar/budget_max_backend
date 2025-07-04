package services

import (
	"net/http"
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/AlsoShantanuBorkar/budget_max/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LoginResponse struct {
	Session     uuid.UUID `json:"session"`
	Refresh     uuid.UUID `json:"refresh"`
	UserID      uuid.UUID `json:"user_id"`
	Requires2FA bool      `json:"requires_2fa"`
	Token       string    `json:"token,omitempty"`
}

type RefreshResponse struct {
	Session uuid.UUID `json:"session"`
	Refresh uuid.UUID `json:"refresh"`
}

func Signup(c *gin.Context, req *models.AuthRequest) *ServiceError {
	// Check if user already exists
	existingUser, err := database.GetUserByEmail(req.Email)
	if err != nil {
		return NewServiceError(http.StatusInternalServerError, "internal server error")
	}

	if existingUser != nil {
		return NewServiceError(http.StatusConflict, "user with this email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return NewServiceError(http.StatusInternalServerError, "failed to process password")
	}

	// Create user
	user := models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	if err := database.CreateUser(&user); err != nil {
		return NewServiceError(http.StatusInternalServerError, "failed to create user")
	}

	return nil
}

func Login(c *gin.Context, req *models.AuthRequest) (*LoginResponse, *ServiceError) {
	// Get user by email
	user, err := database.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		_ = utils.CheckAndTrackLoginAttempts(req.Email)
		return nil, NewServiceError(http.StatusUnauthorized, "invalid email or password")
	}

	// Check password
	if err := utils.CheckPasswordHash(req.Password, user.Password); err != nil {
		_ = utils.CheckAndTrackLoginAttempts(req.Email)
		return nil, NewServiceError(http.StatusUnauthorized, "invalid email or password")
	}

	// Reset login attempts on successful login
	utils.ResetLoginAttempts(req.Email)

	// Check if 2FA is enabled
	if user.TwoFactorEnabled {
		token, err := utils.GenerateJWT(user)
		if err != nil {
			return nil, NewServiceError(http.StatusInternalServerError, "failed to generate token")
		}

		return &LoginResponse{
			Requires2FA: true,
			Token:       token,
		}, nil
	}

	// Create session and refresh token
	session, refresh, err := utils.CreateSessionAndRefreshToken(user, c.ClientIP(), c.Request.UserAgent(), c)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to create session and refresh token")
	}

	return &LoginResponse{
		Session:     session.Token,
		Refresh:     refresh.Token,
		UserID:      user.ID,
		Requires2FA: false,
	}, nil
}

func Logout(c *gin.Context, sessionTokenStr string) *ServiceError {
	if sessionTokenStr == "" {
		return NewServiceError(http.StatusUnauthorized, "unauthorized")
	}

	sessionToken, err := uuid.Parse(sessionTokenStr)
	if err != nil {
		return NewServiceError(http.StatusUnauthorized, "unauthorized")
	}

	session, err := database.GetSessionByToken(sessionToken)
	if err != nil || session == nil {
		return NewServiceError(http.StatusUnauthorized, "invalid session")
	}

	// Revoke session
	if err := database.RevokeSession(session.ID); err != nil {
		return NewServiceError(http.StatusInternalServerError, "failed to revoke session")
	}

	// Revoke all refresh tokens for this session
	if err := database.RevokeRefreshTokensBySessionID(session.ID); err != nil {
		return NewServiceError(http.StatusInternalServerError, "failed to revoke refresh tokens")
	}

	return nil
}

func RefreshToken(c *gin.Context, req *models.RefreshTokensRequest) (*RefreshResponse, *ServiceError) {
	refreshToken, err := uuid.Parse(req.RefreshToken)
	if err != nil {
		return nil, NewServiceError(http.StatusBadRequest, "invalid refresh token")
	}

	// Get refresh token from database
	token, err := database.GetRefreshTokenByToken(refreshToken)
	if err != nil || token == nil || token.Revoked || token.ExpiresAt.Before(time.Now()) {
		return nil, NewServiceError(http.StatusUnauthorized, "invalid or expired refresh token")
	}

	// Revoke the old refresh token
	if err := database.RevokeRefreshToken(token.ID); err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to revoke refresh token")
	}

	// Get user details
	user, err := database.GetUserByID(token.UserID)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to get user details")
	}

	// Create new session and refresh token
	session, refresh, err := utils.CreateSessionAndRefreshToken(user, c.ClientIP(), c.Request.UserAgent(), c)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to create session and refresh token")
	}

	return &RefreshResponse{
		Session: session.Token,
		Refresh: refresh.Token,
	}, nil
}
