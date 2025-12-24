package main

// @title Auth API
// @version 1.0
// @description This API manages orders.
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080

import (
	"context"
	"main/internal/config"
	"main/internal/database"
	"main/internal/handlers"
	"main/internal/logger"
	ginmiddlewares "main/internal/middlewares/gin_middlewares"
	"main/internal/repository"
	"main/internal/services"
	"main/internal/utils"
	"time"

	"github.com/gin-gonic/gin"

	_ "main/cmd/auth/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {


	config.Load()
	isDev:=config.AppCfgs.Server.Env
	dsn:=config.AppCfgs.Database.Postgres
	redis_url:=config.AppCfgs.Database.Redis
	resendApiKey:=config.AppCfgs.Resend.ApiKey

	logger.InitLogger(isDev=="DEV")
	database.Connect(dsn)
	database.InitRedis(redis_url)
	utils.InitJWT()


	utils.InitEmailClient(resendApiKey)
	
	go utils.CleanUpLimits(time.Minute * 5)
	//initializing user repository

	userRepo:=repository.NewRepository(database.DB)
	authService:=services.NewAuthService(userRepo)
	authHandler:=handlers.NewAuthHandler(authService)

	//initializing  routers
	r:=gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth:=r.Group("/auth")
	auth.Use(ginmiddlewares.RateLimit(0.5,10))
	auth.Use(ginmiddlewares.ErrorMiddleware())
	auth.POST("/login",authHandler.Login)
	auth.POST("/verify/email",authHandler.VerifyEmail)
	auth.POST("/refresh",authHandler.Refresh)

	//protected routes
	auth.Use(ginmiddlewares.GinJWTMiddleware())
	{

		auth.Use(ginmiddlewares.RateLimit(2,10))
		auth.GET("/profile",authHandler.Profile)
		auth.GET("/validate",authHandler.Validate)
	}

	ctx:=context.Background()
	go utils.EmailWorker(ctx)
	port:=config.AppCfgs.Server.Port.Auth

	r.Run(":"+port)

}