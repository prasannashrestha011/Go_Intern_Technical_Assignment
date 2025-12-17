package handlers

import (
	"encoding/json"
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



//handlers

func (u *userHandler) GET_ALL_USER(w http.ResponseWriter, r *http.Request) {
	ctx:=r.Context()
	var users []*schema.UserResponseDTO
	users,err:=u.userService.GetUsers(ctx)
	if err!=nil{
		http.Error(w,"User list is empty",400)
		return 
	}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(users)
}

func (u *userHandler) GET_USER(w http.ResponseWriter, r *http.Request) {
	ctx:=r.Context()
	id:=chi.URLParam(r,"id")
	user,err:=u.userService.GetUserByID(ctx,uuid.MustParse(id))
	if err!=nil{
		http.Error(w,err.Error(),404)
		return 
	}
	
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

}

func (u *userHandler) DELETE_USER(w http.ResponseWriter, r *http.Request) {
	ctx:=r.Context()
	id:=chi.URLParam(r,"id")
	err:=u.userService.DeleteUser(ctx,uuid.MustParse(id))
	if err!=nil{
		err:=&utils.CustomError{
			Message: "User not found",
			Status: false,
		}
		http.Error(w,err.Error(),404)
		return
	}
	response:=map[string]string{
		"message":"user deleted with ID: "+id,
	}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(&response)
	
}


func (u *userHandler) REGISTER_USER(w http.ResponseWriter, r *http.Request) {

	ctx:=r.Context()
	var new_user schema.UserCreateDTO 
	if err:=json.NewDecoder(r.Body).Decode(&new_user);err!=nil{
		http.Error(w, err.Error(), 400)
		return
	}
	user_details,err:=u.userService.RegisterUser(ctx,&new_user)
	if err!=nil{
		err:=&utils.CustomError{
			Message:"Email address  exists already",
			Status: false,
		}
		http.Error(w,err.Error(),500)
		return
	}
	w.Header().Set("Content-Type","applicatoin/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&user_details)
}

func (u *userHandler) UPDATE_USER(w http.ResponseWriter, r *http.Request) {

	ctx:=r.Context()
	var userDetails schema.UserUpdateDTO
	id,err:=uuid.Parse(chi.URLParam(r,"id"))
	if err != nil{
		http.Error(w,"Invalid user id",400)
		return
	}
	if err:=json.NewDecoder(r.Body).Decode(&userDetails);err!=nil{
		http.Error(w,"Invalid user details",400)
		return
	}
	 updated_details,err:=u.userService.UpdateUser(ctx,id,&userDetails);
	 if err!=nil{
		http.Error(w,"User Update error: "+err.Error(),500)
		return
	}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(updated_details)
}

