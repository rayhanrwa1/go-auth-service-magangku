package main

import (
	"auth-service/config"
	"auth-service/internal/handler"
	"auth-service/internal/repository"
	"auth-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	db := config.ConnectDB()

	userRepo := repository.NewUserRepository(db)
	refreshRepo := repository.NewRefreshTokenRepository(db)
	resetRepo := repository.NewPasswordResetRepository(db)
	emailService := service.NewEmailService()

	r := gin.Default()

	// TODO : ma
	api := r.Group("/api/v1")
	{
		api.POST("/auth/login", handler.Login(userRepo, refreshRepo))
		api.POST("/auth/logout", handler.Logout(refreshRepo))
		api.POST("/auth/forgot-password", handler.ForgotPassword(userRepo, resetRepo, emailService))
		api.POST("/auth/reset-password", handler.ResetPassword(userRepo, resetRepo))

		register := api.Group("/auth/register")
		{
			register.POST("/user", handler.RegisterByRole(userRepo, refreshRepo, "user"))
			register.POST("/mitra", handler.RegisterByRole(userRepo, refreshRepo, "mitra"))
			register.POST("/admin", handler.RegisterByRole(userRepo, refreshRepo, "admin"))
		}
	}

	r.Run(":8090")
}