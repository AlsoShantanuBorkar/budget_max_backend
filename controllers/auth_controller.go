package controllers

import (
	"fmt"
	"net/http"

	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/AlsoShantanuBorkar/budget_max/services"
	"github.com/AlsoShantanuBorkar/budget_max/utils"
	"github.com/gin-gonic/gin"
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

	serviceErr := services.Signup(c, &req)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
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

	response, serviceErr := services.Login(c, &req)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	if response.Requires2FA {
		c.JSON(http.StatusOK, gin.H{
			"message": "2FA is enabled",
			"data": gin.H{
				"token": response.Token,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data": gin.H{
			"session": response.Session,
			"refresh": response.Refresh,
			"user_id": response.UserID,
		},
	})
}

func Logout(c *gin.Context) {
	sessionTokenStr := c.GetHeader("Authorization")

	serviceErr := services.Logout(c, sessionTokenStr)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

func RefreshToken(c *gin.Context) {
	var req models.RefreshTokensRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(req)
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

	response, serviceErr := services.RefreshToken(c, &req)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"data": gin.H{
			"session": response.Session,
			"refresh": response.Refresh,
			"user_id": response.UserID,
		},
	})
}
