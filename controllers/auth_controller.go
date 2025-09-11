package controllers

import (
	"net/http"

	"github.com/AlsoShantanuBorkar/budget_max/errors"
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/AlsoShantanuBorkar/budget_max/services"
	"github.com/AlsoShantanuBorkar/budget_max/utils"
	"github.com/gin-gonic/gin"
)

type AuthControllerInterface interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	RefreshToken(c *gin.Context)
	Generate2FA(c *gin.Context)
	Verify2FA(c *gin.Context)
	Disable2FA(c *gin.Context)
	LoginWith2FA(c *gin.Context)
}

type AuthController struct {
	service services.AuthServiceInterface
}

func NewAuthController(service services.AuthServiceInterface) *AuthController {
	return &AuthController{service: service}
}

func (ctrl *AuthController) Signup(c *gin.Context) {
	var req models.AuthRequest

       if err := c.ShouldBindJSON(&req); err != nil {
	       appErr := errors.NewBadRequestError("Invalid Request", err, )
	       c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       if err := utils.GetValidator().Struct(req); err != nil {
	       appErr := errors.NewValidationError("Invalid Request", err, )
	       c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       serviceErr := ctrl.service.Signup(c, &req)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr, )
	       c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       c.JSON(http.StatusCreated, gin.H{
	       "message": "User Created Successfully",
       })
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req models.AuthRequest

       if err := c.ShouldBindJSON(&req); err != nil {
	       appErr := errors.NewBadRequestError("Invalid Request", err, )
	       c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       if err := utils.GetValidator().Struct(req); err != nil {
	       appErr := errors.NewValidationError("Invalid Request", err, )
	       c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       response, serviceErr := ctrl.service.Login(c, &req)
       if serviceErr != nil {
	       appErr := errors.NewUnauthorizedError(serviceErr.Message, serviceErr, )
	       c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
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

func (ctrl *AuthController) Logout(c *gin.Context) {
	sessionTokenStr := c.GetHeader("Authorization")

       serviceErr := ctrl.service.Logout(c, sessionTokenStr)
       if serviceErr != nil {
	       appErr := errors.NewUnauthorizedError(serviceErr.Message, serviceErr, )
	       c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       c.JSON(http.StatusOK, gin.H{
	       "message": "Logged out successfully",
       })
}

func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	var req models.RefreshTokensRequest

       if err := c.ShouldBindJSON(&req); err != nil {
	       appErr := errors.NewBadRequestError("Invalid Request", err, )
	       c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       if err := utils.GetValidator().Struct(req); err != nil {
	       appErr := errors.NewValidationError("Invalid Request", err, )
	       c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       response, serviceErr := ctrl.service.RefreshToken(c, &req)
       if serviceErr != nil {
	       appErr := errors.NewUnauthorizedError(serviceErr.Message, serviceErr, )
	       c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
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

func (ctrl *AuthController) Generate2FA(c *gin.Context) {
       userId, ok := utils.ParseUserID(c)
       if !ok {
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil, )
	       c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       response, serviceErr := ctrl.service.Generate2FA(c, userId)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr, )
	       c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       c.JSON(http.StatusOK, gin.H{
	       "message": "2FA key generated successfully",
	       "data":    response,
       })
}

func (ctrl *AuthController) Verify2FA(c *gin.Context) {
       userId, ok := utils.ParseUserID(c)
       if !ok {
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil, )
	       c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       var request models.TwoFactorVerifyRequest
       if err := c.ShouldBindJSON(&request); err != nil {
	       appErr := errors.NewBadRequestError("Invalid Request", err, )
	       c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       if err := utils.GetValidator().Struct(request); err != nil {
	       appErr := errors.NewValidationError("Invalid Request", err, )
	       c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       serviceErr := ctrl.service.Verify2FA(c, &request, userId)
       if serviceErr != nil {
	       appErr := errors.NewUnauthorizedError(serviceErr.Message, serviceErr, )
	       c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       c.JSON(http.StatusOK, gin.H{
	       "message": "2FA verified successfully",
       })
}

func (ctrl *AuthController) Disable2FA(c *gin.Context) {
       userId, ok := utils.ParseUserID(c)
       if !ok {
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil, )
	       c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       serviceErr := ctrl.service.Disable2FA(c, userId)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr, )
	       c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       c.JSON(http.StatusOK, gin.H{
	       "message": "2FA disabled successfully",
       })
}

func (ctrl *AuthController) LoginWith2FA(c *gin.Context) {
       var request models.TwoFactorLoginRequest
       if err := c.ShouldBindJSON(&request); err != nil {
	       appErr := errors.NewBadRequestError("Invalid Request", err, )
	       c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       if err := utils.GetValidator().Struct(request); err != nil {
	       appErr := errors.NewValidationError("Invalid Request", err, )
	       c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       response, serviceErr := ctrl.service.LoginWith2FA(c, &request)
       if serviceErr != nil {
	       appErr := errors.NewUnauthorizedError(serviceErr.Message, serviceErr, )
	       c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
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
