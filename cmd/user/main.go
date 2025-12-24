package main

// @title User API
// @version 1.0
// @description This API manages orders.
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8081
import (
	"context"
	"log"
	_ "main/cmd/user/docs"
	"main/internal/config"
	"main/internal/database"
	"main/internal/handlers"
	"main/internal/logger"
	chimiddlewares "main/internal/middlewares/chi_middlewares"
	"main/internal/repository"
	"main/internal/services"
	"main/internal/utils"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	config.Load()
	isDev:=config.AppCfgs.Server.Env
	dsn:=config.AppCfgs.Database.Postgres
	redis_url:=config.AppCfgs.Database.Redis
	resendApiKey:=config.AppCfgs.Resend.ApiKey

	logger.InitLogger(isDev=="DEV")

	utils.InitJWT()
	database.Connect(dsn)
	database.InitRedis(redis_url)

	
	utils.InitJWT()


	utils.InitEmailClient(resendApiKey)



	r := chi.NewRouter()

	repo:=repository.NewRepository(database.DB)
	userService:=services.NewUserService(repo)
	userHandlers:=handlers.NewUserHandler(userService)
	
	go utils.CleanUpLimits(time.Minute * 5)

	r.Use(chimiddlewares.LoggerMiddleware)
	r.Use(chimiddlewares.ErrorMiddleware)

	r.Use(chimiddlewares.RateLimit(0.5,10))
	r.Post("/create",userHandlers.REGISTER_USER)
	r.Route("/users",func(r chi.Router) {
		r.Use(chimiddlewares.JWTAuthMiddleware)
		r.Get("/",userHandlers.GET_ALL_USER)
		r.Get("/{id}",userHandlers.GET_USER)
		r.Put("/update/{id}",userHandlers.UPDATE_USER)
		r.Delete("/{id}",userHandlers.DELETE_USER)
	})
	
	port:=config.AppCfgs.Server.Port.User
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:"+port+"/swagger/doc.json"), // URL pointing to generated swagger.json
	))

	ctx:=context.Background()
	go utils.EmailWorker(ctx)
	log.Println("SERVER listening on PORT: "+port)
	http.ListenAndServe(":"+port,r)
}
