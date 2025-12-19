package main

// @title Auth API
// @version 1.0
// @description This API manages orders.
// @termsOfService http://example.com/terms/
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

import (
	"main/internal/config"
	"main/internal/database"
	"main/internal/handlers"
	"main/internal/logger"
	ginmiddlewares "main/internal/middlewares/gin_middlewares"
	"main/internal/repository"
	"main/internal/services"
	"main/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	_ "main/cmd/auth/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {


	config.Load()
	isDev:=config.AppCfgs.Server.Env
	dsn:=config.AppCfgs.Database.Url

	logger.InitLogger(isDev=="DEV")
	err:=database.Connect(dsn)
	if err!=nil{
		logger.Log.Error("Database connection error: ",zap.Error(err))
	}

	utils.InitJWT()
	//initializing user repository

	userRepo:=repository.NewRepository(database.DB)
	authService:=services.NewAuthService(userRepo)
	authHandler:=handlers.NewAuthHandler(authService)

	//initializing  routers
	r:=gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth:=r.Group("/api/v1/auth")
	auth.Use(ginmiddlewares.ErrorMiddleware())
	auth.POST("/login",authHandler.Login)
	auth.POST("/refresh",authHandler.Refresh)

	//protected routes
	auth.Use(ginmiddlewares.GinJWTMiddleware())
	{
		auth.GET("/profile",authHandler.Profile)
		auth.GET("/validate",authHandler.Validate)
	}

	port:=config.AppCfgs.Server.Port.Auth

	r.Run(":"+port)

}