package services

import (
	"net/http"

	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/AlsoShantanuBorkar/budget_max/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
)

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

func Generate2FA(c *gin.Context, userId uuid.UUID) (*TwoFAGenerateResponse, *ServiceError) {
	user, err := database.GetUserByID(userId)
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

	err = database.UpdateUser(userId, map[string]any{
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

func Verify2FA(c *gin.Context, req *models.TwoFactorVerifyRequest, userId uuid.UUID) *ServiceError {
	user, err := database.GetUserByID(userId)
	if err != nil {
		return NewServiceError(http.StatusInternalServerError, "failed to get user")
	}

	valid := totp.Validate(req.Code, user.TwoFactorSecret)
	if !valid {
		return NewServiceError(http.StatusUnauthorized, "invalid 2FA code")
	}

	user.TwoFactorEnabled = true

	err = database.UpdateUser(userId, map[string]any{
		"two_factor_enabled": user.TwoFactorEnabled,
	})

	if err != nil {
		return NewServiceError(http.StatusInternalServerError, "failed to update user")
	}

	return nil
}

func Disable2FA(c *gin.Context, userId uuid.UUID) *ServiceError {
	user, err := database.GetUserByID(userId)
	if err != nil {
		return NewServiceError(http.StatusInternalServerError, "failed to get user")
	}

	if !user.TwoFactorEnabled {
		return NewServiceError(http.StatusBadRequest, "2FA is not enabled")
	}

	user.TwoFactorEnabled = false

	err = database.UpdateUser(userId, map[string]any{
		"two_factor_enabled": user.TwoFactorEnabled,
		"two_factor_secret":  "",
	})

	if err != nil {
		return NewServiceError(http.StatusInternalServerError, "failed to update user")
	}

	return nil
}

func LoginWith2FA(c *gin.Context, req *models.TwoFactorLoginRequest) (*TwoFALoginResponse, *ServiceError) {
	token, claims, err := utils.VerifyJWT(req.Token)

	if err != nil || !token || !claims.Is2FA {
		return nil, NewServiceError(http.StatusUnauthorized, "token is invalid or expired")
	}

	if claims.Email != req.Email {
		return nil, NewServiceError(http.StatusUnauthorized, "invalid email")
	}

	user, err := database.GetUserByEmail(claims.Email)
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
	session, refresh, err := utils.CreateSessionAndRefreshToken(user, c.ClientIP(), c.Request.UserAgent(), c)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to create session and refresh token")
	}

	return &TwoFALoginResponse{
		Session: session.Token,
		Refresh: refresh.Token,
		UserID:  user.ID,
	}, nil
}
