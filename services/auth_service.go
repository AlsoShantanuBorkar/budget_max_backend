package services

import (
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/config"
	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/errors"
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/AlsoShantanuBorkar/budget_max/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
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
	config                      *config.AppConfig
	redisClient                 *redis.Client
}

func NewAuthService(userDBService database.UserDatabaseServiceInterface, sessionDBService database.SessionDatabaseServiceInterface, refreshTokenDBService database.RefreshTokenDatabaseServiceInterface, config *config.AppConfig, redisClient *redis.Client) *AuthService {
	return &AuthService{
		userDatabaseService:         userDBService,
		sessionDatabaseService:      sessionDBService,
		refreshTokenDatabaseService: refreshTokenDBService,
		config:                      config,
		redisClient:                 redisClient,
	}
}

func (s *AuthService) Signup(c *gin.Context, req *models.AuthRequest) *ServiceError {
	// Check if user already exists
       existingUser, err := s.userDatabaseService.GetUserByEmail(req.Email)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil
       }

       if existingUser != nil {
	       appErr := errors.NewConflictError("user with this email already exists", nil, )
	       c.Error(appErr)
	       return nil
       }

       // Hash password
       hashedPassword, err := utils.HashPassword(req.Password)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil
       }

       // Create user
       user := models.User{
	       ID:        uuid.New(),
	       Email:     req.Email,
	       Password:  hashedPassword,
	       CreatedAt: time.Now(),
       }

       if err := s.userDatabaseService.CreateUser(&user); err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil
       }

       return nil
}

func (s *AuthService) Login(c *gin.Context, req *models.AuthRequest) (*LoginResponse, *ServiceError) {
	// Get user by email
       user, err := s.userDatabaseService.GetUserByEmail(req.Email)
       if err != nil || user == nil {
	       if trackErr := utils.CheckAndTrackLoginAttempts(req.Email, s.redisClient, c.Request.Context()); trackErr != nil {
		       appErr := errors.NewTooManyRequestsError(trackErr.Error(), trackErr, )
		       c.Error(appErr)
		       return nil, nil
	       }
	       appErr := errors.NewUnauthorizedError("invalid email or password", err, )
	       c.Error(appErr)
	       return nil, nil
       }

       // Check password
       if err := utils.CheckPasswordHash(req.Password, user.Password); err != nil {
	       if trackErr := utils.CheckAndTrackLoginAttempts(req.Email, s.redisClient, c.Request.Context()); trackErr != nil {
		       // If the account is now locked, return that specific message.
		       appErr := errors.NewTooManyRequestsError(trackErr.Error(), trackErr, )
		       c.Error(appErr)
		       return nil, nil
	       }
	       appErr := errors.NewUnauthorizedError("invalid email or password", err, )
	       c.Error(appErr)
	       return nil, nil
       }

       // Reset login attempts on successful login
	utils.ResetLoginAttempts(req.Email, s.redisClient, c.Request.Context())

       // Check if 2FA is enabled
       if user.TwoFactorEnabled {
	       token, err := utils.GenerateJWT(user, s.config)
	       if err != nil {
		       appErr := errors.NewInternalError(err, )
		       c.Error(appErr)
		       return nil, nil
	       }

	       return &LoginResponse{
		       Requires2FA: true,
		       Token:       token,
	       }, nil
       }

       // Create session and refresh token
       session, refresh, err := utils.CreateSessionAndRefreshToken(user, c.ClientIP(), c.Request.UserAgent(), s.sessionDatabaseService, s.refreshTokenDatabaseService, c)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, nil
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
	       appErr := errors.NewUnauthorizedError("unauthorized", nil, )
	       c.Error(appErr)
	       return nil
       }

       sessionToken, err := uuid.Parse(sessionTokenStr)
       if err != nil {
	       appErr := errors.NewUnauthorizedError("unauthorized", err, )
	       c.Error(appErr)
	       return nil
       }

       session, err := s.sessionDatabaseService.GetSessionByToken(sessionToken)
       if err != nil || session == nil {
	       appErr := errors.NewUnauthorizedError("invalid session", err, )
	       c.Error(appErr)
	       return nil
       }

       // Revoke session
       if err := s.sessionDatabaseService.RevokeSession(session.ID); err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil
       }

       // Revoke all refresh tokens for this session
       if err := s.refreshTokenDatabaseService.RevokeRefreshTokensBySessionID(session.ID); err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil
       }

       return nil
}

func (s *AuthService) RefreshToken(c *gin.Context, req *models.RefreshTokensRequest) (*RefreshResponse, *ServiceError) {
       refreshToken, err := uuid.Parse(req.RefreshToken)
       if err != nil {
	       appErr := errors.NewBadRequestError("invalid refresh token", err, )
	       c.Error(appErr)
	       return nil, nil
       }

       // Get refresh token from database
       token, err := s.refreshTokenDatabaseService.GetRefreshTokenByToken(refreshToken)
       if err != nil || token == nil || token.Revoked || token.ExpiresAt.Before(time.Now()) {
	       appErr := errors.NewUnauthorizedError("invalid or expired refresh token", err, )
	       c.Error(appErr)
	       return nil, nil
       }

       // Revoke the old refresh token
       if err := s.refreshTokenDatabaseService.RevokeRefreshToken(token.ID); err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, nil
       }

       // Get user details
       user, err := s.userDatabaseService.GetUserByID(token.UserID)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, nil
       }

       // Create new session and refresh token
       session, refresh, err := utils.CreateSessionAndRefreshToken(user, c.ClientIP(), c.Request.UserAgent(), s.sessionDatabaseService, s.refreshTokenDatabaseService, c)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, nil
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
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, nil
       }

       if user.TwoFactorEnabled {
	       appErr := errors.NewBadRequestError("2FA is already enabled", nil, )
	       c.Error(appErr)
	       return nil, nil
       }

       secret, err := totp.Generate(totp.GenerateOpts{
	       Issuer:      "BudgetMax",
	       AccountName: user.Email,
       })

       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, nil
       }

       user.TwoFactorSecret = secret.Secret()

       err = s.userDatabaseService.UpdateUser(userId, map[string]any{
	       "two_factor_secret":  user.TwoFactorSecret,
	       "two_factor_enabled": user.TwoFactorEnabled,
       })

       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, nil
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
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil
       }

       valid := totp.Validate(req.Code, user.TwoFactorSecret)
       if !valid {
	       appErr := errors.NewUnauthorizedError("invalid 2FA code", nil, )
	       c.Error(appErr)
	       return nil
       }

       user.TwoFactorEnabled = true

       err = s.userDatabaseService.UpdateUser(userId, map[string]any{
	       "two_factor_enabled": user.TwoFactorEnabled,
       })

       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil
       }

       return nil
}

func (s *AuthService) Disable2FA(c *gin.Context, userId uuid.UUID) *ServiceError {
       user, err := s.userDatabaseService.GetUserByID(userId)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil
       }

       if !user.TwoFactorEnabled {
	       appErr := errors.NewBadRequestError("2FA is not enabled", nil, )
	       c.Error(appErr)
	       return nil
       }

       user.TwoFactorEnabled = false

       err = s.userDatabaseService.UpdateUser(userId, map[string]any{
	       "two_factor_enabled": user.TwoFactorEnabled,
	       "two_factor_secret":  "",
       })

       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil
       }

       return nil
}

func (s *AuthService) LoginWith2FA(c *gin.Context, req *models.TwoFactorLoginRequest) (*TwoFALoginResponse, *ServiceError) {
       token, claims, err := utils.VerifyJWT(req.Token, s.config)

       if err != nil || !token || !claims.Is2FA {
	       appErr := errors.NewUnauthorizedError("token is invalid", err, )
	       c.Error(appErr)
	       return nil, nil
       }

       if claims.RegisteredClaims.ExpiresAt.Before(time.Now()) {
	       appErr := errors.NewUnauthorizedError("token is expired", nil, )
	       c.Error(appErr)
	       return nil, nil
       }

       if claims.Email != req.Email {
	       appErr := errors.NewUnauthorizedError("invalid email", nil, )
	       c.Error(appErr)
	       return nil, nil
       }

       user, err := s.userDatabaseService.GetUserByEmail(claims.Email)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, nil
       }

       if !user.TwoFactorEnabled {
	       appErr := errors.NewBadRequestError("2FA is not enabled", nil, )
	       c.Error(appErr)
	       return nil, nil
       }

       valid := totp.Validate(req.Code, user.TwoFactorSecret)
       if !valid {
	       appErr := errors.NewUnauthorizedError("invalid 2FA code", nil, )
	       c.Error(appErr)
	       return nil, nil
       }

       // Create session and refresh token
       session, refresh, err := utils.CreateSessionAndRefreshToken(user, c.ClientIP(), c.Request.UserAgent(), s.sessionDatabaseService, s.refreshTokenDatabaseService, c)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, nil
       }

       return &TwoFALoginResponse{
	       Session: session.Token,
	       Refresh: refresh.Token,
	       UserID:  user.ID,
       }, nil
}
