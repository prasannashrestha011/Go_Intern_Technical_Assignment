package handlers

import (
	"main/internal/logger"
	"main/internal/schema"
	"main/internal/services"
	"main/internal/utils"
	"net/http"
	"time"

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
		_=ctx.Error(utils.ErrBadRequest)
		return
	}

	logger.Log.Info("Login Attempted",zap.String("Email",userCreds.Email),zap.Time("logged_at",time.Now()))
	authData,err:=a.authService.Login(ctx,userCreds)

	if err!=nil{
		logger.Log.Info("Login Failed: ",zap.String("Email",userCreds.Email))
		_=ctx.Error(utils.ErrInvalidCredentials)
		return

	}

	accessToken,refreshToken,err:=utils.GenerateTokens(authData.ID)
	if err!=nil{
		_=ctx.Error(utils.ErrTokenInvalid)
		return
	}

	logger.Log.Info("Login Success: ",zap.String("Email",userCreds.Email))
	ctx.SetCookie("refresh_token",refreshToken,7*24*3600,"/","",false,true)
	ctx.Header("Authorization","Bearer "+accessToken)
	ctx.JSON(http.StatusOK,gin.H{
		"message":"login successful",
	})
}

func (a *authHandler) Refresh(ctx *gin.Context) {
	refreshToken,err:=ctx.Cookie("refresh_token")
	if err!=nil{
		_=ctx.Error(utils.ErrTokenMissing)
		return
	}
	token,err:=utils.ValidateJWT(refreshToken)
	if err!=nil{
		_=ctx.Error(utils.ErrTokenInvalid)
		return
	}
	userIDStr,err:=utils.GenerateUserIDFromToken(token)
	if err!=nil{
		_=ctx.Error(utils.ErrTokenInvalid)
	}

	newAccessToken,_,err:=utils.GenerateTokens(uuid.MustParse(userIDStr))
	if err!=nil{
		_=ctx.Error(utils.ErrTokenGenFailure)
	}
	logger.Log.Info("Refresh token generated: ",zap.String("userID",userIDStr),zap.Time("time",time.Now()))
	ctx.JSON(http.StatusOK,gin.H{
		"accessToken":"Bearer "+newAccessToken,
	})
}

func (a *authHandler) Validate(ctx *gin.Context) {
	userID,exists:=ctx.Get("userID")
	if !exists{
		_=ctx.Error(utils.ErrTokenMissing)
		return
	}
	logger.Log.Info("Attempted Token validation: ",zap.String("userID",userID.(string)),zap.Time("time",time.Now()))
	ctx.JSON(http.StatusOK,gin.H{
		"message":"Token is valid",
		"userID":userID,
	})
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
	logger.Log.Info("Logged user",zap.String("userID",userID.(string)),zap.Time("time",time.Now()))
	ctx.JSON(http.StatusOK,gin.H{
		"message":"You are an authorized user",
		"userID":userID,
	})
}

