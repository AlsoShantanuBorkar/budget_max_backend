package controllers

import (
	"net/http"
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/AlsoShantanuBorkar/budget_max/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Signup(c *gin.Context) {

	var req models.AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
		})
		return
	}
	if err := utils.GetValidator().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
		})
		return
	}

	existingUser, err := database.GetUserByEmail(req.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{
			"message": "User with this email already exists",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to process password",
		})
		return
	}

	user := models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	err = database.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create User",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User Created Successfully",
	})

}
func Login(c *gin.Context) {
	var req models.AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
		})
		return
	}

	if err := utils.GetValidator().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
		})
		return
	}

	user, err := database.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid Email or Password",
		})
		return
	}

	if err := utils.CheckPasswordHash(req.Password, user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid Email or Password",
		})
		return
	}

	if user.TwoFactorEnabled {
		token, err := utils.GenerateJWT(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to generate token",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "2FA is enabled",
			"data": gin.H{
				"token": token,
			},
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

func Logout(c *gin.Context) {
	sessionTokenStr := c.GetHeader("Authorization")
	if sessionTokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	sessionToken, err := uuid.Parse(sessionTokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
	}

	session, err := database.GetSessionByToken(sessionToken)
	if err != nil || session == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid session",
		})
		return
	}

	if err := database.RevokeSession(session.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to revoke session",
		})
		return
	}

	if err := database.RevokeRefreshTokensBySessionID(session.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to revoke refresh tokens",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

func RefreshToken(c *gin.Context) {
	var req models.RefreshTokensRequest
	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
		})
		return
	}
	if err := utils.GetValidator().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
		})
		return
	}

	refreshToken, err := uuid.Parse(req.RefreshToken)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Refresh Token",
		})
		return
	}

	token, err := database.GetRefreshTokenByToken(refreshToken)
	if err != nil || token == nil || token.Revoked || token.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid or expired refresh token",
		})
		return
	}

	if err := database.RevokeRefreshToken(token.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to revoke refresh token",
		})
		return
	}

	user, err := database.GetUserByID(token.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get user details",
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
		"message": "Token refreshed successfully",
		"data": gin.H{
			"session": session.Token,
			"refresh": refresh.Token,
		},
	})
}
