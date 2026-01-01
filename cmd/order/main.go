package main

// @title Orders API
// @version 1.0
// @description This API manages orders.
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8082
import (
	"log"
	"main/internal/config"
	"main/internal/config/protoc"
	"main/internal/database"
	"main/internal/handlers"
	"main/internal/logger"
	chimiddlewares "main/internal/middlewares/chi_middlewares"
	"main/internal/repository"
	"main/internal/services"
	"main/internal/utils"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	userClient := protoc.NewUserMicroServiceClient(conn)

	config.Load()
	isDev := config.AppCfgs.Server.Env
	dsn := config.AppCfgs.Database.Postgres
	logger.InitLogger(isDev == "DEV")

	utils.InitJWT()
	database.Connect(dsn)

	utils.InitOrderWorker()

	repo := repository.NewOrderRepository(database.DB)
	service := services.NewOrderService(repo)
	handler := handlers.NewOrderHandler(service, userClient)

	r := mux.NewRouter()

	r.Use(chimiddlewares.RateLimit(0.5, 10))
	r.Use(chimiddlewares.ErrorMiddleware)
	r.Use(chimiddlewares.JWTAuthMiddleware)

	go utils.CleanUpLimits(time.Minute * 5)

	r.HandleFunc("/orders", handler.GetALLOrders).Methods("GET")
	r.HandleFunc("/orders", handler.CreateOrder).Methods("POST")

	r.HandleFunc("/orders/{id}", handler.GetOrder).Methods("GET")
	r.HandleFunc("/orders/{id}", handler.UpdateOrderDetails).Methods("PUT")
	r.HandleFunc("/users/{id}/orders", handler.GetUserOrders).Methods("GET")

	http.ListenAndServe(":"+config.AppCfgs.Server.Port.Order, r)
}
