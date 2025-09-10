package services

import (
	"net/http"
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/AlsoShantanuBorkar/budget_max/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
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
	UserID  uuid.UUID `json:"user_id"`
}

type TwoFAGenerateResponse struct {
	Secret     string `json:"secret"`
	OTPAuthURL string `json:"otp_auth_url"`
	Issuer     string `json:"issuer"`
	Email      string `json:"email"`
}

type TwoFALoginResponse struct {
	Session uuid.UUID `json:"session"`
	Refresh uuid.UUID `json:"refresh"`
	UserID  uuid.UUID `json:"user_id"`
}

type AuthServiceInterface interface {
	Signup(c *gin.Context, req *models.AuthRequest) *ServiceError
	Login(c *gin.Context, req *models.AuthRequest) (*LoginResponse, *ServiceError)
	Logout(c *gin.Context, sessionTokenStr string) *ServiceError
	RefreshToken(c *gin.Context, req *models.RefreshTokensRequest) (*RefreshResponse, *ServiceError)
	Generate2FA(c *gin.Context, userId uuid.UUID) (*TwoFAGenerateResponse, *ServiceError)
	Verify2FA(c *gin.Context, req *models.TwoFactorVerifyRequest, userId uuid.UUID) *ServiceError
	Disable2FA(c *gin.Context, userId uuid.UUID) *ServiceError
	LoginWith2FA(c *gin.Context, req *models.TwoFactorLoginRequest) (*TwoFALoginResponse, *ServiceError)
}

type AuthService struct {
	userDatabaseService         database.UserDatabaseServiceInterface
	sessionDatabaseService      database.SessionDatabaseServiceInterface
	refreshTokenDatabaseService database.RefreshTokenDatabaseServiceInterface
}

func NewAuthService(userDBService database.UserDatabaseServiceInterface, sessionDBService database.SessionDatabaseServiceInterface, refreshTokenDBService database.RefreshTokenDatabaseServiceInterface) *AuthService {
	return &AuthService{
		userDatabaseService:         userDBService,
		sessionDatabaseService:      sessionDBService,
		refreshTokenDatabaseService: refreshTokenDBService,
	}
}

func (s *AuthService) Signup(c *gin.Context, req *models.AuthRequest) *ServiceError {
	// Check if user already exists
	existingUser, err := s.userDatabaseService.GetUserByEmail(req.Email)
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

	if err := s.userDatabaseService.CreateUser(&user); err != nil {
		return NewServiceError(http.StatusInternalServerError, "failed to create user")
	}

	return nil
}

func (s *AuthService) Login(c *gin.Context, req *models.AuthRequest) (*LoginResponse, *ServiceError) {
	// Get user by email
	user, err := s.userDatabaseService.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		if trackErr := utils.CheckAndTrackLoginAttempts(req.Email); trackErr != nil {
			return nil, NewServiceError(http.StatusTooManyRequests, trackErr.Error())
		}
		return nil, NewServiceError(http.StatusUnauthorized, "invalid email or password")
	}

	// Check password
	if err := utils.CheckPasswordHash(req.Password, user.Password); err != nil {
		if trackErr := utils.CheckAndTrackLoginAttempts(req.Email); trackErr != nil {
			// If the account is now locked, return that specific message.
			return nil, NewServiceError(http.StatusTooManyRequests, trackErr.Error())
		}
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
	session, refresh, err := utils.CreateSessionAndRefreshToken(user, c.ClientIP(), c.Request.UserAgent(), s.sessionDatabaseService, s.refreshTokenDatabaseService, c)
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

func (s *AuthService) Logout(c *gin.Context, sessionTokenStr string) *ServiceError {
	if sessionTokenStr == "" {
		return NewServiceError(http.StatusUnauthorized, "unauthorized")
	}

	sessionToken, err := uuid.Parse(sessionTokenStr)
	if err != nil {
		return NewServiceError(http.StatusUnauthorized, "unauthorized")
	}

	session, err := s.sessionDatabaseService.GetSessionByToken(sessionToken)
	if err != nil || session == nil {
		return NewServiceError(http.StatusUnauthorized, "invalid session")
	}

	// Revoke session
	if err := s.sessionDatabaseService.RevokeSession(session.ID); err != nil {
		return NewServiceError(http.StatusInternalServerError, "failed to revoke session")
	}

	// Revoke all refresh tokens for this session
	if err := s.refreshTokenDatabaseService.RevokeRefreshTokensBySessionID(session.ID); err != nil {
		return NewServiceError(http.StatusInternalServerError, "failed to revoke refresh tokens")
	}

	return nil
}

func (s *AuthService) RefreshToken(c *gin.Context, req *models.RefreshTokensRequest) (*RefreshResponse, *ServiceError) {
	refreshToken, err := uuid.Parse(req.RefreshToken)
	if err != nil {
		return nil, NewServiceError(http.StatusBadRequest, "invalid refresh token")
	}

	// Get refresh token from database
	token, err := s.refreshTokenDatabaseService.GetRefreshTokenByToken(refreshToken)
	if err != nil || token == nil || token.Revoked || token.ExpiresAt.Before(time.Now()) {
		return nil, NewServiceError(http.StatusUnauthorized, "invalid or expired refresh token")
	}

	// Revoke the old refresh token
	if err := s.refreshTokenDatabaseService.RevokeRefreshToken(token.ID); err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to revoke refresh token")
	}

	// Get user details
	user, err := s.userDatabaseService.GetUserByID(token.UserID)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to get user details")
	}

	// Create new session and refresh token
	session, refresh, err := utils.CreateSessionAndRefreshToken(user, c.ClientIP(), c.Request.UserAgent(), s.sessionDatabaseService, s.refreshTokenDatabaseService, c)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to create session and refresh token")
	}

	return &RefreshResponse{
		Session: session.Token,
		Refresh: refresh.Token,
		UserID:  user.ID,
	}, nil
}

func (s *AuthService) Generate2FA(c *gin.Context, userId uuid.UUID) (*TwoFAGenerateResponse, *ServiceError) {
	user, err := s.userDatabaseService.GetUserByID(userId)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to get user")
	}

	if user.TwoFactorEnabled {
		return nil, NewServiceError(http.StatusBadRequest, "2FA is already enabled")
	}

	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "BudgetMax",
		AccountName: user.Email,
	})

	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to generate key")
	}

	user.TwoFactorSecret = secret.Secret()

	err = s.userDatabaseService.UpdateUser(userId, map[string]any{
		"two_factor_secret":  user.TwoFactorSecret,
		"two_factor_enabled": user.TwoFactorEnabled,
	})

	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to update user")
	}

	return &TwoFAGenerateResponse{
		Secret:     user.TwoFactorSecret,
		OTPAuthURL: string(secret.URL()),
		Issuer:     secret.Issuer(),
		Email:      user.Email,
	}, nil
}

func (s *AuthService) Verify2FA(c *gin.Context, req *models.TwoFactorVerifyRequest, userId uuid.UUID) *ServiceError {
	user, err := s.userDatabaseService.GetUserByID(userId)
	if err != nil {
		return NewServiceError(http.StatusInternalServerError, "failed to get user")
	}

	valid := totp.Validate(req.Code, user.TwoFactorSecret)
	if !valid {
		return NewServiceError(http.StatusUnauthorized, "invalid 2FA code")
	}

	user.TwoFactorEnabled = true

	err = s.userDatabaseService.UpdateUser(userId, map[string]any{
		"two_factor_enabled": user.TwoFactorEnabled,
	})

	if err != nil {
		return NewServiceError(http.StatusInternalServerError, "failed to update user")
	}

	return nil
}

func (s *AuthService) Disable2FA(c *gin.Context, userId uuid.UUID) *ServiceError {
	user, err := s.userDatabaseService.GetUserByID(userId)
	if err != nil {
		return NewServiceError(http.StatusInternalServerError, "failed to get user")
	}

	if !user.TwoFactorEnabled {
		return NewServiceError(http.StatusBadRequest, "2FA is not enabled")
	}

	user.TwoFactorEnabled = false

	err = s.userDatabaseService.UpdateUser(userId, map[string]any{
		"two_factor_enabled": user.TwoFactorEnabled,
		"two_factor_secret":  "",
	})

	if err != nil {
		return NewServiceError(http.StatusInternalServerError, "failed to update user")
	}

	return nil
}

func (s *AuthService) LoginWith2FA(c *gin.Context, req *models.TwoFactorLoginRequest) (*TwoFALoginResponse, *ServiceError) {
	token, claims, err := utils.VerifyJWT(req.Token)

	if err != nil || !token || !claims.Is2FA {
		return nil, NewServiceError(http.StatusUnauthorized, "token is invalid")
	}

	if claims.RegisteredClaims.ExpiresAt.Before(time.Now()) {
		return nil, NewServiceError(http.StatusUnauthorized, "token is expired")
	}

	if claims.Email != req.Email {
		return nil, NewServiceError(http.StatusUnauthorized, "invalid email")
	}

	user, err := s.userDatabaseService.GetUserByEmail(claims.Email)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to get user")
	}

	if !user.TwoFactorEnabled {
		return nil, NewServiceError(http.StatusBadRequest, "2FA is not enabled")
	}

	valid := totp.Validate(req.Code, user.TwoFactorSecret)
	if !valid {
		return nil, NewServiceError(http.StatusUnauthorized, "invalid 2FA code")
	}

	// Create session and refresh token
	session, refresh, err := utils.CreateSessionAndRefreshToken(user, c.ClientIP(), c.Request.UserAgent(), s.sessionDatabaseService, s.refreshTokenDatabaseService, c)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to create session and refresh token")
	}

	return &TwoFALoginResponse{
		Session: session.Token,
		Refresh: refresh.Token,
		UserID:  user.ID,
	}, nil
}
