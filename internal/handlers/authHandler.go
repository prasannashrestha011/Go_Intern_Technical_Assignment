package handlers

import (
	"main/internal/logger"
	"main/internal/schema"
	"main/internal/services"
	"main/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthHandler interface {
	Login(ctx *gin.Context)
	Refresh(ctx *gin.Context)
	Validate(ctx *gin.Context)
	//test handlers
	Profile(ctx *gin.Context)
}

type authHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) AuthHandler {
	return &authHandler{authService: authService}
}


func (a *authHandler) Login(ctx *gin.Context) {
	var userCreds *schema.UserLoginDTO

	if err:=ctx.BindJSON(&userCreds);err!=nil{
		logger.Log.Error("Request Body error: ",zap.Error(err))
		ctx.JSON(http.StatusBadRequest,gin.H{
			"error":"Invalid request body",
		})
	}

	authData,err:=a.authService.Login(ctx,userCreds)
	if err!=nil{
		ctx.JSON(http.StatusForbidden,gin.H{
			"error":"Invalid credentials",
		})
	}
	accessToken,refreshToken,err:=utils.GenerateTokens(authData.ID)
	if err!=nil{
		ctx.JSON(http.StatusForbidden,gin.H{
			"error":"Failed to generate authentication token, please try again",
		})
	}
	ctx.SetCookie("refresh_token",refreshToken,7*24*3600,"/","",false,true)
	ctx.Header("Authorization","Bearer "+accessToken)
	ctx.JSON(http.StatusOK,gin.H{
		"message":"login successful",
	})
}

func (a *authHandler) Refresh(ctx *gin.Context) {
	refreshToken,err:=ctx.Cookie("refresh_token")
	if err!=nil{
		ctx.JSON(http.StatusUnauthorized,gin.H{
			"error":"Refresh token missing !!",
		})
		return
	}
	token,err:=utils.ValidateJWT(refreshToken)
	if err!=nil{
		ctx.JSON(http.StatusUnauthorized,gin.H{
			"error":"Invalid refresh token, please login again",
		})
		return
	}
	userIDStr,err:=utils.GenerateUserIDFromToken(token)
	if err!=nil{
		ctx.JSON(http.StatusUnauthorized,gin.H{
			"error":"Invalid refresh token, please login again",
		})
	}

	newAccessToken,_,err:=utils.GenerateTokens(uuid.MustParse(userIDStr))
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"error":"Failed to generate new access token",
		})
	}
	ctx.JSON(http.StatusOK,gin.H{
		"accessToken":"Bearer "+newAccessToken,
	})
}

func (a *authHandler) Validate(ctx *gin.Context) {
	panic("unimplemented")
}

/*
	Test handlers to test jwt route protection
*/

func (a * authHandler) Profile(ctx *gin.Context){
	userID,exists:=ctx.Get("userID")
	if !exists{
		ctx.JSON(http.StatusUnauthorized,gin.H{
			"error":"User Id not found in the context",
		})
		return
	}
	logger.Log.Info("Logged user",zap.String("userID",userID.(string)))
	ctx.JSON(http.StatusOK,gin.H{
		"message":"You are an authorized user",
		"userID":userID,
	})
}

