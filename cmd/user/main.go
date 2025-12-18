package main

import (
	"log"
	"main/internal/config"
	"main/internal/database"
	"main/internal/handlers"
	"main/internal/logger"
	chimiddlewares "main/internal/middlewares/chi_middlewares"
	"main/internal/repository"
	"main/internal/services"
	"main/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	config.Load()
	isDev:=config.AppCfgs.Server.Env
	dsn:=config.AppCfgs.Database.Url
	logger.InitLogger(isDev=="DEV")

	utils.InitJWT()
	database.Connect(dsn)
	r := chi.NewRouter()

	repo:=repository.NewRepository(database.DB)
	service:=services.NewUserService(repo)
	userHandlers:=handlers.NewUserHandler(service)
	
	r.Use(chimiddlewares.LoggerMiddleware)
	r.Route("/users",func(r chi.Router) {
		r.Use(chimiddlewares.JWTAuthMiddleware)
		r.Get("/",userHandlers.GET_ALL_USER)
		r.Get("/{id}",userHandlers.GET_ALL_USER)
		r.Post("/create",userHandlers.REGISTER_USER)
		r.Put("/update/{id}",userHandlers.UPDATE_USER)
		r.Delete("/{id}",userHandlers.DELETE_USER)
	})
	

	port:=config.AppCfgs.Server.Port.User

	log.Println("SERVER listening on PORT: "+port)
	http.ListenAndServe(":"+port,r)
}
