package handlers

import (
	"encoding/json"
	"main/internal/config/protoc"
	chimiddlewares "main/internal/middlewares/chi_middlewares"
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
	UpdateOrderDetails(w http.ResponseWriter, r *http.Request)
}

type orderHandler struct {
	orderService services.OrderService
	userClient   protoc.UserMicroServiceClient
}

func NewOrderHandler(orderService services.OrderService, userClient protoc.UserMicroServiceClient) OrderHandler {
	return &orderHandler{
		orderService: orderService,
		userClient:   userClient,
	}
}

// CreateOrder godoc
// @Summary      Create a new order
// @Description  Creates a new order for a user with the provided order data
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        order  body      schema.CreateOrder  true  "Order payload"
// @Success      201    {object}  schema.SuccessResponseSchema "Order created successfully"
// @Failure      400    {object}  schema.ErrorResponseSchema "Invalid order details"
// @Failure      500    {object}  schema.ErrorResponseSchema "Unknow error occured in the server"
// @Router       /orders [post]
func (o *orderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var newOrder schema.CreateOrder
	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		chimiddlewares.SetError(w, utils.NewAppError(http.StatusBadRequest, "BAD_REQUEST", "Invalid Email address or password", nil))
		return
	}

	order, err := o.orderService.CreateOrder(ctx, &newOrder)

	if err != nil {

		chimiddlewares.SetError(w, utils.NewAppError(http.StatusInternalServerError, "INTERNAL_SERVER_ERR", "Failed to create the new order", nil))
		return
	}
	response := schema.SuccessResponse(order, "Order created")
	utils.JsonResponseWriter(w, http.StatusCreated, response)
}

// GetALLOrders godoc
// @Summary      List all orders
// @Description  Returns a list of all orders in the system
// @Tags         orders
// @Produce      json
// @Success      200  {object}  schema.SuccessResponseSchema "Order lists from the database"
// @Failure      404  {object}  schema.ErrorResponseSchema	"Order list is empty in the database"
// @Router       /orders [get]
func (o *orderHandler) GetALLOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page, pageSize := utils.ParsePaginationValues(r, 1, 10)

	orderList, err := o.orderService.GetOrders(ctx, page, pageSize)
	if err != nil {
		chimiddlewares.SetError(w, utils.NewAppError(http.StatusNotFound, "ENTITY_NOT_FOUND", "No order details in the database", nil))
		return
	}

	resp := schema.SuccessResponse(orderList, "Order list")
	utils.JsonResponseWriter(w, http.StatusOK, resp)
}

// GetOrder godoc
// @Summary      Get order by ID
// @Description  Returns details of a specific order by its ID
// @Tags         orders
// @Produce      json
// @Param        id   path      string  true  "Order ID"
// @Success      200  {object}  schema.SuccessResponseSchema "Order details associated with orderID: "
// @Failure      400  {object}  schema.ErrorResponseSchema "Order Id is missing"
// @Failure      404  {object}  schema.ErrorResponseSchema "Order details not found"
// @Router       /orders/{id} [get]
func (o *orderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	idStr, ok := params["id"]
	if !ok {
		chimiddlewares.SetError(w, utils.NewAppError(http.StatusBadRequest, "BAD_REQUEST", "Order ID is missing", nil))
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		chimiddlewares.SetError(w, utils.NewAppError(http.StatusBadRequest, "BAD_REQUEST", "Invalid Order ID, it must be in UUID format", nil))
		return
	}
	order, err := o.orderService.GetOrder(ctx, id)
	if err != nil {
		chimiddlewares.SetError(w, utils.NewAppError(http.StatusNotFound, "ENTITY_NOT_FOUND", "Order details not found", nil))
		return
	}
	customer_details, err := o.userClient.GetUser(ctx, &protoc.UserMicroRequest{UserId: order.UserID.String()})
	if err != nil {
		chimiddlewares.SetError(w, utils.NewAppError(http.StatusNotFound, "ENTITY_NOT_FOUND", "Customer details not found", nil))
		return
	}
	final_details := map[string]interface{}{
		"order":    order,
		"customer": customer_details,
	}
	resp := schema.SuccessResponse(final_details, "Order details")
	utils.JsonResponseWriter(w, http.StatusOK, resp)
}

// GetUserOrders godoc
// @Summary      Get all orders for a user
// @Description  Returns a list of orders associated with a specific user ID
// @Tags         orders
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  schema.SuccessResponseSchema "Order details associated with userID: "
// @Failure      400  {object}  schema.ErrorResponseSchema	"UserId is invalid or missing!!"
// @Failure      404  {object}  schema.ErrorResponseSchema	"Order details didnot found with userID: "
// @Router       /users/{id}/orders [get]
func (o *orderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	idStr, ok := params["id"]
	page, pageSize := utils.ParsePaginationValues(r, 1, 10)
	if !ok {
		chimiddlewares.SetError(w, utils.NewAppError(http.StatusBadRequest, "BAD_REQUEST", "Invalid User ID, it must be in UUID format", nil))
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		chimiddlewares.SetError(w, utils.NewAppError(http.StatusBadRequest, "BAD_REQUEST", "Invalid Order ID, it must be in UUID format", nil))
		return
	}
	orderList, err := o.orderService.GetUserOrders(ctx, id, page, pageSize)
	if err != nil {
		chimiddlewares.SetError(w, utils.NewAppError(http.StatusNotFound, "ENTITY_NOT_FOUND", "Order details not found with userID: "+idStr, nil))
		return
	}
	resp := schema.SuccessResponse(orderList, "Order details associated with UserID: "+id.String())
	utils.JsonResponseWriter(w, http.StatusOK, resp)
}

// UpdateOrderDetails godoc
// @Summary      Update an existing order
// @Description  Updates the details of a specific order by its ID
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        id     path      string              true  "Order ID"
// @Param        order  body      schema.OrderUpdate  true  "Updated order payload"
// @Success      200    {object}  schema.SuccessResponseSchema "Order updated successfully"
// @Failure      400    {object}  schema.ErrorResponseSchema   "Invalid Order ID or order payload"
// @Failure      404    {object}  schema.ErrorResponseSchema   "Order not found"
// @Failure      500    {object}  schema.ErrorResponseSchema   "Internal server error"
// @Router       /orders/{id} [put]
func (o *orderHandler) UpdateOrderDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	idStr, ok := params["id"]
	if !ok {
		chimiddlewares.SetError(w, utils.NewAppError(http.StatusBadRequest, "BAD_REQUEST", "Invalid User ID, it must be in UUID format", nil))
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		chimiddlewares.SetError(w, utils.NewAppError(http.StatusBadRequest, "BAD_REQUEST", "Invalid Order ID, it must be in UUID format", nil))
		return
	}
	var orderDto *schema.OrderUpdate
	if err := json.NewDecoder(r.Body).Decode(&orderDto); err != nil {
		resp := utils.NewAppError(http.StatusBadRequest, "BAD_REQUEST", "Invalid order details", nil)
		chimiddlewares.SetError(w, resp)
		return
	}
	updatedOrder, err := o.orderService.UpdateOrderDetails(ctx, id, orderDto)
	if err != nil {
		resp := utils.NewAppError(http.StatusBadRequest, "BAD_REQUEST", "Failed to update order details with orderID: "+idStr, nil)
		chimiddlewares.SetError(w, resp)
		return
	}
	utils.JsonResponseWriter(w, http.StatusOK, updatedOrder)
}
