package main

import (
	"io"
	"log"
	"net/http"
	"os"

	ginmiddlewares "main/internal/middlewares/gin_middlewares"
	"main/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	prometheus.MustRegister(utils.HttpRequestsTotal)
	prometheus.MustRegister(utils.ActiveSessionsGuage)
	prometheus.MustRegister(utils.HttpRequestDuration)
	logFile := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	}

	// Write to both stdout and file
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	gin.DefaultWriter = multiWriter
}

func main() {
	r := gin.Default()

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/health", gin.WrapH(promhttp.Handler()))
	r.POST("/login", ginmiddlewares.MetricMiddleware(), func(ctx *gin.Context) {
		utils.ActiveSessionsGuage.Inc()
		log.Println("User logged in ")
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Session logged In",
		})
	})
	r.POST("/logout", ginmiddlewares.MetricMiddleware(), func(ctx *gin.Context) {
		utils.ActiveSessionsGuage.Dec()

		log.Println("User logged out ")
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Session logged out",
		})
	})
	r.GET("/request/distribution", ginmiddlewares.MetricMiddleware(), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "request distribution",
		})
	})

	log.Println("Server running on port: 2112")
	if err := r.Run(":2112"); err != nil {
		log.Println(err.Error())
	}
}
