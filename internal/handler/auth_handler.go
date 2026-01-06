package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"auth-service/internal/domain"
	"auth-service/internal/middleware"
	"auth-service/internal/repository"
	"auth-service/internal/service"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type RegisterRequest struct {
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Password      string `json:"password" binding:"required,min=8"`
	TermsAccepted bool   `json:"terms_accepted"`
}

func RegisterByRole(
	userRepo *repository.UserRepository,
	refreshRepo *repository.RefreshTokenRepository,
	role string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Message: "invalid request",
				Data:    nil,
			})
			return
		}

		if !req.TermsAccepted {
			c.JSON(http.StatusBadRequest, Response{
				Message: "terms must be accepted",
				Data:    nil,
			})
			return
		}

		hashedPassword, err := service.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Message: "failed to hash password",
				Data:    nil,
			})
			return
		}

		user := domain.User{
			ID:            uuid.NewString(),
			Name:          req.Name,
			Email:         req.Email,
			Password:      hashedPassword,
			UserableType:  role,
			TermsAccepted: true,
		}

		if err := userRepo.Create(context.Background(), &user); err != nil {
			c.JSON(http.StatusConflict, Response{
				Message: "user already exists",
				Data:    nil,
			})
			return
		}

		accessToken, err := middleware.GenerateAccessToken(user.ID, user.Name, user.UserableType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Message: "failed to generate access token",
				Data:    nil,
			})
			return
		}

		refreshToken, exp, err := middleware.GenerateRefreshToken(user.ID, user.Name, user.UserableType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Message: "failed to generate refresh token",
				Data:    nil,
			})
			return
		}

		_ = refreshRepo.Create(
			c.Request.Context(),
			user.ID,
			refreshToken,
			exp,
		)

		c.JSON(http.StatusCreated, Response{
			Message: "success",
			Data: gin.H{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			},
		})
	}
}

func Login(
	userRepo *repository.UserRepository,
	refreshRepo *repository.RefreshTokenRepository,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Message: "invalid request",
				Data:    nil,
			})
			return
		}

		user, err := userRepo.FindByEmail(c.Request.Context(), req.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, Response{
				Message: "invalid email or password",
				Data:    nil,
			})
			return
		}

		if err := service.CheckPassword(user.Password, req.Password); err != nil {
			c.JSON(http.StatusUnauthorized, Response{
				Message: "invalid email or password",
				Data:    nil,
			})
			return
		}

		accessToken, err := middleware.GenerateAccessToken(user.ID, user.Name, user.UserableType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Message: "failed to generate access token",
				Data:    nil,
			})
			return
		}

		refreshToken, exp, err := middleware.GenerateRefreshToken(user.ID, user.Name, user.UserableType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Message: "failed to generate refresh token",
				Data:    nil,
			})
			return
		}

		_ = refreshRepo.Create(
			c.Request.Context(),
			user.ID,
			refreshToken,
			exp,
		)

		c.JSON(http.StatusOK, Response{
			Message: "success",
			Data: gin.H{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			},
		})
	}
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

func Logout(refreshRepo *repository.RefreshTokenRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Message: "invalid request",
				Data:    nil,
			})
			return
		}

		claims, err := middleware.ParseRefreshToken(req.RefreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, Response{
				Message: "invalid token",
				Data:    nil,
			})
			return
		}

		userID := claims["user_id"].(string)

		if err := refreshRepo.Delete(c.Request.Context(), userID, req.RefreshToken); err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Message: "failed to logout",
				Data:    nil,
			})
			return
		}

		c.JSON(http.StatusOK, Response{
			Message: "logout successful",
			Data:    nil,
		})
	}
}

func ForgotPassword(
	userRepo *repository.UserRepository, 
	resetRepo *repository.PasswordResetRepository,
	emailService *service.EmailService,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ForgotPasswordRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Message: "invalid request",
				Data:    nil,
			})
			return
		}

		user, err := userRepo.FindByEmail(c.Request.Context(), req.Email)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				Message: "if email exists, reset link will be sent",
				Data:    nil,
			})
			return
		}

		token, exp, err := service.GenerateResetToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Message: "failed to generate reset token",
				Data:    nil,
			})
			return
		}

		if err := resetRepo.Create(c.Request.Context(), user.ID, token, exp); err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Message: "failed to save reset token",
				Data:    nil,
			})
			return
		}

		go emailService.SendResetPasswordEmail(user.Email, token)

		c.JSON(http.StatusOK, Response{
			Message: "if email exists, reset link will be sent",
			Data:    nil,
		})
	}
}

func ResetPassword(userRepo *repository.UserRepository, resetRepo *repository.PasswordResetRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ResetPasswordRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Message: "invalid request",
				Data:    nil,
			})
			return
		}

		userID, err := resetRepo.Verify(c.Request.Context(), req.Token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, Response{
				Message: "invalid or expired token",
				Data:    nil,
			})
			return
		}

		hashedPassword, err := service.HashPassword(req.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Message: "failed to hash password",
				Data:    nil,
			})
			return
		}

		if err := userRepo.UpdatePassword(c.Request.Context(), userID, hashedPassword); err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Message: "failed to update password",
				Data:    nil,
			})
			return
		}

		if err := resetRepo.Delete(c.Request.Context(), req.Token); err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Message: "failed to delete reset token",
				Data:    nil,
			})
			return
		}

		c.JSON(http.StatusOK, Response{
			Message: "password reset successful",
			Data:    nil,
		})
	}
}