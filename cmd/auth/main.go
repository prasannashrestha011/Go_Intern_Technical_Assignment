package main

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
	auth:=r.Group("/auth")
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