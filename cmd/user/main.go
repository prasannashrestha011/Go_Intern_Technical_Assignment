package user

import (
	"log"
	"main/internal/database"
	"main/internal/handlers"
	"main/internal/repository"
	"main/internal/services"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func USER_CMD() {
	err:=godotenv.Load()
	if err!=nil{
		log.Println("ENV error: ",err.Error())
	}

	dsn:=os.Getenv("DB_URL")
	database.Connect(dsn)
	r := chi.NewRouter()

	repo:=repository.NewRepository(database.DB)
	service:=services.NewUserService(repo)
	userHandlers:=handlers.NewUserHandler(service)
	
	r.Get("/users/{id}",userHandlers.GET_USER)
	r.Get("/users",userHandlers.GET_ALL_USER)
	r.Post("/users/create",userHandlers.REGISTER_USER)
	r.Put("/users/update/{id}",userHandlers.UPDATE_USER)
	r.Delete("/users/{id}",userHandlers.DELETE_USER)

	port:=os.Getenv("PORT")

	log.Println("SERVER listening on PORT: "+port)
	http.ListenAndServe(":"+port,r)
}
