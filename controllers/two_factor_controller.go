package controllers

import (
	"net/http"

	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/AlsoShantanuBorkar/budget_max/services"
	"github.com/AlsoShantanuBorkar/budget_max/utils"
	"github.com/gin-gonic/gin"
)

func Generate2FA(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	response, serviceErr := services.Generate2FA(c, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "2FA key generated successfully",
		"data":    response,
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

	if err := utils.GetValidator().Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	serviceErr := services.Verify2FA(c, &request, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
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

	serviceErr := services.Disable2FA(c, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
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

	if err := utils.GetValidator().Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	response, serviceErr := services.LoginWith2FA(c, &request)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "2FA login successful",
		"data": gin.H{
			"session": response.Session,
			"refresh": response.Refresh,
			"user_id": response.UserID,
		},
	})
}
