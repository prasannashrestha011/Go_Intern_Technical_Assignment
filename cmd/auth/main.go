package auth

import (
	"log"
	"main/internal/database"
	"main/internal/handlers"
	"main/internal/logger"
	"main/internal/middlewares"
	"main/internal/repository"
	"main/internal/services"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func AUTH_CMD() {
	err := godotenv.Load()
	if err!=nil{
		log.Println("Failed to laod the env: ",err.Error())
	}

	isDev:=os.Getenv("ENV")
	logger.InitLogger(isDev=="DEV")

	dsn:=os.Getenv("DB_URL")
	err=database.Connect(dsn)
	if err!=nil{
		logger.Log.Error("Database connection error: ",zap.Error(err))
	}

	//initializing user repository

	userRepo:=repository.NewRepository(database.DB)
	authService:=services.NewAuthService(userRepo)
	authHandler:=handlers.NewAuthHandler(authService)

	//initializing  routers
	r:=gin.Default()
	auth:=r.Group("/auth")
	auth.POST("/login",authHandler.Login)
	auth.Use(middlewares.GinJWTMiddleware())
	{
		auth.GET("/profile",authHandler.Profile)
	}

	port:=os.Getenv("PORT")

	r.Run(":"+port)

}