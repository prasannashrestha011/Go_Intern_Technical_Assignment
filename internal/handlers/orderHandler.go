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

// CreateOrder godoc
// @Summary      Create a new order
// @Description  Creates a new order for a user with the provided order data
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        order  body      schema.CreateOrder  true  "Order payload"
// @Success      201    {object}  schema.Response
// @Failure      400    {object}  schema.Response
// @Failure      500    {object}  schema.Response
// @Router       /orders [post]
func (o *orderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx:=r.Context()
	var newOrder schema.CreateOrder
	if err:=json.NewDecoder(r.Body).Decode(&newOrder);err!=nil{
		resp:=utils.ErrBadRequest.WithDetails("Invalid email address or password")
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if err:=o.orderService.CreateOrder(ctx,&newOrder);err!=nil{
		resp := schema.ErrorResponse("INTERNAL_SERVER_ERR","Failed to create the order",err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	response:=schema.SuccessResponse(nil,"Order created")
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&response)
}

// GetALLOrders godoc
// @Summary      List all orders
// @Description  Returns a list of all orders in the system
// @Tags         orders
// @Produce      json
// @Success      200  {object}  schema.Response
// @Failure      404  {object}  schema.Response
// @Router       /orders [get]
func (o *orderHandler) GetALLOrders(w http.ResponseWriter, r *http.Request) {
	ctx:=r.Context()
	orderList,err:=o.orderService.GetOrders(ctx)
	if err!=nil{
		resp:=schema.ErrorResponse("NOT_FOUND","Order list is empty","")
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&resp)
		return
	}

	resp:=schema.SuccessResponse(orderList,"Order list")
	w.Header().Set("Content-Type","application/json")	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}


// GetOrder godoc
// @Summary      Get order by ID
// @Description  Returns details of a specific order by its ID
// @Tags         orders
// @Produce      json
// @Param        id   path      string  true  "Order ID"
// @Success      200  {object}  schema.Response
// @Failure      400  {object}  schema.Response
// @Failure      404  {object}  schema.Response
// @Router       /orders/{id} [get]
func (o *orderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx:=r.Context()
	params:=mux.Vars(r)
	idStr,ok:=params["id"]
	if !ok{
		resp:=schema.ErrorResponse("BAD_REQUEST","Order ID is missing","")
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&resp)
		return
	}
	id,err:=uuid.Parse(idStr)
	if err!=nil{
		resp:=schema.ErrorResponse("BAD_REQUEST","Invalid order ID","")
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&resp)
		return
	}
	order,err:=o.orderService.GetOrder(ctx,id)
	if err!=nil{
		resp:=schema.ErrorResponse("NOT_FOUND","Order details not found","")
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&resp)
		return
	}
	resp:=schema.SuccessResponse(order,"Order details")
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(&resp)
}


// GetUserOrders godoc
// @Summary      Get all orders for a user
// @Description  Returns a list of orders associated with a specific user ID
// @Tags         orders
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  schema.Response
// @Failure      400  {object}  schema.Response
// @Failure      404  {object}  schema.Response
// @Router       /users/{id}/orders [get]
func (o *orderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	ctx:=r.Context()
	params:=mux.Vars(r)
	idStr,ok:=params["id"]
	if !ok{
		resp:=schema.ErrorResponse("BAD_REQUEST","User ID is missing","")
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&resp)
		return
	}
	id,err:=uuid.Parse(idStr)
	if err!=nil{
		resp:=schema.ErrorResponse("BAD_REQUEST","Invalid order ID","")
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&resp)
	}
	orderList,err:=o.orderService.GetUserOrders(ctx,id)
	if err!=nil{
		resp:=schema.ErrorResponse("NOT_FOUND","Orders not found with userID","")
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&resp)
		return
	}
	resp:=schema.SuccessResponse(orderList,"Order details associated with UserID: "+id.String())
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(&resp)
}