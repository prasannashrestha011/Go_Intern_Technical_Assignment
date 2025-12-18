package handlers

import (
	"encoding/json"
	"main/internal/schema"
	"main/internal/services"
	"main/internal/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type OrderHandler interface {
	GetOrder(w http.ResponseWriter, r *http.Request)
	GetUserOrders(w http.ResponseWriter, r *http.Request)
	GetALLOrders(w http.ResponseWriter, r *http.Request)
	CreateOrder(w http.ResponseWriter, r *http.Request)
}

type orderHandler struct {
	orderService services.OrderService
}


func NewOrderHandler(orderService services.OrderService) OrderHandler {
	return &orderHandler{
		orderService: orderService,
	}
}

func (o *orderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx:=r.Context()
	var newOrder schema.CreateOrder
	if err:=json.NewDecoder(r.Body).Decode(&newOrder);err!=nil{
		err:=utils.ErrBadRequest
		http.Error(w,err.Message,err.Code)
		return
	}

	err:=o.orderService.CreateOrder(ctx,&newOrder)
	if err!=nil{
		http.Error(w,"Failed to create the error",500)
		return
	}
	response:=&schema.ResponseSchema{
		Message: "Order created with userID "+newOrder.UserID.String(),
		Code: http.StatusCreated,
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(response.Code)
	json.NewEncoder(w).Encode(&response)
}

func (o *orderHandler) GetALLOrders(w http.ResponseWriter, r *http.Request) {
	ctx:=r.Context()
	orderList,err:=o.orderService.GetOrders(ctx)
	if err!=nil{
		err:=utils.AppError{
			Message: "No order details found",
			Code: http.StatusNotFound,
			Err: nil,
		}
		http.Error(w,err.Message,err.Code)
		return
	}
	
	w.Header().Set("Content-Type","application/json")	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orderList)
}

func (o *orderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx:=r.Context()
	params:=mux.Vars(r)
	idStr,ok:=params["id"]
	if !ok{
		http.Error(w,"Order Id is missing!!",http.StatusBadRequest)
		return
	}
	id,err:=uuid.Parse(idStr)
	if err!=nil{
		http.Error(w,"Invalid order ID, uuid format didnt matched",http.StatusBadRequest)
		return
	}
	order,err:=o.orderService.GetOrder(ctx,id)
	if err!=nil{
		http.Error(w,"Order not found",http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(order)
}

func (o *orderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	ctx:=r.Context()
	params:=mux.Vars(r)
	idStr,ok:=params["id"]
	if !ok{
		http.Error(w,"User ID is missing !!",http.StatusBadRequest)
		return
	}
	id,err:=uuid.Parse(idStr)
	if err!=nil{
		http.Error(w,"Invalid user ID, uuid format didnt matched",http.StatusBadRequest)
		return
	}
	orderList,err:=o.orderService.GetUserOrders(ctx,id)
	if err!=nil{
		http.Error(w,"Orders not found with userID"+id.String(),http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(orderList)
}