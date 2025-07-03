package controllers

import (
	"fmt"
	"net/http"

	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/AlsoShantanuBorkar/budget_max/utils"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
)

func Generate2FA(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	user, err := database.GetUserByID(userId)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get user",
		})
		return
	}

	if user.TwoFactorEnabled {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "2FA is already enabled",
		})
		return
	}

	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "BudgetMax",
		AccountName: user.Email,
	})

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate key",
		})
		return
	}

	user.TwoFactorSecret = secret.Secret()

	err = database.UpdateUser(userId, map[string]any{
		"two_factor_secret":  user.TwoFactorSecret,
		"two_factor_enabled": user.TwoFactorEnabled,
	})

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update user",
		})
		return
	}
	url := string(secret.URL())
	issuer := secret.Issuer()

	c.JSON(http.StatusOK, gin.H{
		"message": "2FA key generated successfully",
		"data": gin.H{
			"secret":       user.TwoFactorSecret,
			"otp_auth_url": url,
			"issuer":       issuer,
			"email":        user.Email,
		},
	})
}

func Verify2FA(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	var request models.TwoFactorVerifyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, err := database.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get user",
		})
		return
	}

	valid := totp.Validate(request.Code, user.TwoFactorSecret)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid 2FA code",
		})
		return
	}

	user.TwoFactorEnabled = true

	err = database.UpdateUser(userId, map[string]any{
		"two_factor_enabled": user.TwoFactorEnabled,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "2FA verified successfully",
	})
}

func Disable2FA(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	user, err := database.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get user",
		})
		return
	}

	if !user.TwoFactorEnabled {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "2FA is not enabled",
		})
		return
	}

	user.TwoFactorEnabled = false

	err = database.UpdateUser(userId, map[string]any{
		"two_factor_enabled": user.TwoFactorEnabled,
		"two_factor_secret":  "",
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "2FA disabled successfully",
	})
}

func LoginWith2FA(c *gin.Context) {

	var request models.TwoFactorLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	token, claims, err := utils.VerifyJWT(request.Token)

	if err != nil || !token || !claims.Is2FA {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token is invalid or expired",
		})
		return
	}

	if claims.Email != request.Email {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid email",
		})
		return
	}

	user, err := database.GetUserByEmail(claims.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get user",
		})
		return
	}

	if !user.TwoFactorEnabled {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "2FA is not enabled",
		})
		return
	}

	valid := totp.Validate(request.Code, user.TwoFactorSecret)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid 2FA code",
		})
		return
	}

	session, refresh, err := utils.CreateSessionAndRefreshToken(user, c.ClientIP(), c.Request.UserAgent(), c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create session and refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data": gin.H{
			"session": session.Token,
			"refresh": refresh.Token,
			"user_id": user.ID,
		},
	})
}
