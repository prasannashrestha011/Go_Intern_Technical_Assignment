package handlers

import (
	"encoding/json"
	"fmt"
	"main/internal/config"
	chimiddlewares "main/internal/middlewares/chi_middlewares"
	"main/internal/schema"
	"main/internal/services"
	"main/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserService interface {
	GET_USER(w http.ResponseWriter, r *http.Request)
	GET_ALL_USER(w http.ResponseWriter, r *http.Request)
	REGISTER_USER(w http.ResponseWriter, r *http.Request)
	UPDATE_USER(w http.ResponseWriter, r *http.Request)
	DELETE_USER(w http.ResponseWriter, r *http.Request)
}
type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) UserService {
	return &userHandler{userService:userService}

}

// GET_ALL_USER godoc
// @Summary      List all users
// @Description  Returns a list of all users from the database
// @Tags         users
// @Produce      json
// @Success      200  {object}  schema.SuccessResponseSchema
// @Failure      404  {object}  schema.ErrorResponseSchema
// @Router       /users [get]
func (u *userHandler) GET_ALL_USER(w http.ResponseWriter, r *http.Request) {
	ctx:=r.Context()
	var users []*schema.UserResponseDTO
	users,err:=u.userService.GetUsers(ctx)
	if err!=nil{
		appErr:=utils.NewAppError(http.StatusNotFound,"ENTITY_NOT_FOUND","No users found in the system",nil)
		chimiddlewares.SetError(w,appErr)
		return 
	}
	resp:=schema.SuccessResponse(users,"List of all users from the database")
	utils.JsonResponseWriter(w,http.StatusOK,resp)
}


// GET_USER godoc
// @Summary      Get user by ID
// @Description  Returns user details for a specific user ID
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  schema.SuccessResponseSchema
// @Failure      404  {object}  schema.ErrorResponseSchema
// @Router       /users/{id} [get]
func (u *userHandler) GET_USER(w http.ResponseWriter, r *http.Request) {
	ctx:=r.Context()
	id:=chi.URLParam(r,"id")
	user,err:=u.userService.GetUserByID(ctx,uuid.MustParse(id))
	if err!=nil{
		appErr:=utils.NewAppError(http.StatusNotFound,"ENTITY_NOT_FOUND","User details not found",nil)
		chimiddlewares.SetError(w, appErr)
		return
	}
	
	resp:=schema.SuccessResponse(user,"User details associated with user ID: "+id)
	utils.JsonResponseWriter(w,http.StatusOK,resp)
}


// DELETE_USER godoc
// @Summary      Delete a user
// @Description  Deletes a user by their user ID
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  schema.SuccessResponseSchema
// @Failure      400  {object}  schema.ErrorResponseSchema
// @Failure      500  {object}  schema.ErrorResponseSchema
// @Router       /users/{id} [delete]
func (u *userHandler) DELETE_USER(w http.ResponseWriter, r *http.Request) {
	ctx:=r.Context()
	idStr:=chi.URLParam(r,"id")
	id,err:=uuid.Parse(idStr)
	if err!=nil{
		appErr:=utils.NewAppError(http.StatusInternalServerError,"BAD_REQUEST","User ID didnt matched the UUID format",nil)
		chimiddlewares.SetError(w,appErr)
	}
	err=u.userService.DeleteUser(ctx,id)
	if err!=nil{
		appErr:=utils.NewAppError(http.StatusInternalServerError,"INTERNAL_SERVER_ERR","Failed to delete user details with userID: "+idStr,err)
		chimiddlewares.SetError(w,appErr)
		return
	}
	resp:=schema.SuccessResponse(nil,"User details delete  with user ID: "+idStr)
	utils.JsonResponseWriter(w,http.StatusOK,resp)
	
}

// REGISTER_USER godoc
// @Summary      Register a new user
// @Description  Creates a new user with the provided details
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      schema.UserCreateDTO  true  "User registration payload"
// @Success      200   {object}  schema.SuccessResponseSchema
// @Failure      400   {object}  schema.ErrorResponseSchema
// @Failure      500   {object}  schema.ErrorResponseSchema
// @Router       /users/create [post]
func (u *userHandler) REGISTER_USER(w http.ResponseWriter, r *http.Request) {

	ctx:=r.Context()
	var new_user schema.UserCreateDTO 
	if err:=json.NewDecoder(r.Body).Decode(&new_user);err!=nil{
		http.Error(w, err.Error(), 400)
		return
	}
	user_details,err:=u.userService.RegisterUser(ctx,&new_user)
	if err!=nil{
		appErr:=utils.NewAppError(http.StatusInternalServerError,"INTERNAL_SERVER_ERR","Email address already exists",nil)
		chimiddlewares.SetError(w,appErr)
		return
	}

	code:=utils.GenerateVerificationCode()
	err=utils.StoreVerificationCode(ctx,user_details.Email,code)
	if err!=nil{
		resp:=utils.NewAppError(http.StatusInternalServerError,"INTERNAL_SERVER_ERR","Failed to generate verification code, please try again later",nil)
		chimiddlewares.SetError(w,resp)
		return
	}
	job_data:=&utils.EmailJob{
		To: user_details.Email,
		From: config.AppCfgs.Resend.AppDomain,
		Subj: "Account verification code",
		Body: fmt.Sprintf(`
			<h2>Email Verification</h2>
			<p>Hi %s,</p>
			<p>Thank you for registering. Please use the verification code below to complete your registration:</p>
			<h3 style="color: #2F80ED;">%d</h3>
			<p>This code will expire in 3 minutes.</p>
			<p>If you did not request this, please ignore this email.</p>
			<hr>
			<p style="font-size: 12px; color: #888;">Need help? Contact our support team at support@example.com</p>
		`, user_details.Name, code),
	}
	utils.EnqueueEmail(ctx,job_data)
	resp:=schema.SuccessResponse(nil,"Verification code has been sent to "+user_details.Email)
	utils.JsonResponseWriter(w,http.StatusOK,resp)
}


// UPDATE_USER godoc
// @Summary      Update user details
// @Description  Updates user details for a given user ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string              true  "User ID"
// @Param        user  body      schema.UserUpdateDTO  true  "Updated user data"
// @Success      200   {object}  schema.SuccessResponseSchema
// @Failure      400   {object}  schema.ErrorResponseSchema
// @Failure      500   {object}  schema.ErrorResponseSchema
// @Router       /users/update/{id} [put]
func (u *userHandler) UPDATE_USER(w http.ResponseWriter, r *http.Request) {

	ctx:=r.Context()
	var userDetails schema.UserUpdateDTO
	id,err:=uuid.Parse(chi.URLParam(r,"id"))
	if err != nil{
		appErr:=utils.NewAppError(http.StatusBadRequest,"BAD_REQUEST","User ID is missing",nil)
		chimiddlewares.SetError(w,appErr)
		return
	}
	if err:=json.NewDecoder(r.Body).Decode(&userDetails);err!=nil{
		appErr:=utils.NewAppError(http.StatusBadRequest,"BAD_REQUEST","Invalid request body",nil)
		chimiddlewares.SetError(w,appErr)
		return
	}
	 updated_details,err:=u.userService.UpdateUser(ctx,id,&userDetails);
	 if err!=nil{
		appErr:=utils.NewAppError(http.StatusInternalServerError,"INTERNAL_SERVER_ERR","Failed to update user details, please try again",nil)
		chimiddlewares.SetError(w,appErr)
		return
	}
	resp:=schema.SuccessResponse(updated_details,"User details  updated with email "+updated_details.Email)
	utils.JsonResponseWriter(w,http.StatusOK,resp)
}

