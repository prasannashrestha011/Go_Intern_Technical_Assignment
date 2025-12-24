package handlers

import (
	"main/internal/config"
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
	VerifyEmail(ctx *gin.Context)
	//test handlers
	Profile(ctx *gin.Context)
}

type authHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) AuthHandler {
	return &authHandler{authService: authService }
}

// Login godoc
//@Summary Login function for the app
//@Description The request body will contain fields (email and password) as credentials. The given credentials will be used to authenticate user. If authentication is successful, access and refresh token(jwt) will be issued to the user with claims(userID).
//@Tags auth
//@Accept json
//@Produce json
//@Param userCreds body schema.UserLoginDTO true "User credentials"
//@Success 200 {object} schema.SuccessResponseSchema "Login successful"
//@Error 400 {object} schema.ErrorResponseSchema "Bad request"
//@Error 401 {object} schema.ErrorResponseSchema "Invalid credentials"
//@Router /auth/login [post]
func (a *authHandler) Login(ctx *gin.Context) {
	var userCreds *schema.UserLoginDTO

	if err:=ctx.BindJSON(&userCreds);err!=nil{
		_=ctx.Error(utils.ErrBadRequest)
		return
	}

	logger.Log.Info("Login Attempted",zap.String("Email",userCreds.Email),zap.Time("logged_at",time.Now()))

	authData,err:=a.authService.Login(ctx,userCreds)
	if err!=nil{
		_=ctx.Error(err)
		return

	}

	accessToken,refreshToken,err:=utils.GenerateTokens(authData.ID)
	if err!=nil{
		_=ctx.Error(utils.NewAppError(401,"TOKEN_GEN_FAILURE","Failed to generate authorization tokens",nil))
		return
	}

	logger.Log.Info("Login Success: ",zap.String("Email",userCreds.Email))

	ctx.SetCookie("refresh_token",refreshToken,7*24*3600,"/","",false,true)
	responseData := map[string]string{
		"accessToken": accessToken,
	}

	ctx.JSON(http.StatusOK,schema.SuccessResponse(responseData,"Login successful"))
}
// Refresh Auth Token godoc
//@Summary Refresh access token through refresh token
//@Description This handler is resposible for refreshing the session state of access token. The access token persist for 15m after its issue. Once the token is expired, refresh token which is issued during login action, is used in order to revise the session state of access token.
//Tags auth
//@Accept json
//@Produce json
//Param refresh_token cookie string true "Refresh token cookie"
//@Success 200 {object} schema.SuccessResponseSchema "New access token"
//@Failure 401 {object} schema.ErrorResponseSchema "Invalid refresh token"
//@Router /auth/refresh [post]
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
		_=ctx.Error(utils.NewAppError(500,"TOKEN_GEN_FAILTURE","Failed to refresh access token",nil))
		return
	}
	logger.Log.Info("Refresh token generated: ",zap.String("userID",userIDStr),zap.Time("time",time.Now()))
	responseData := map[string]string{
		"accessToken": newAccessToken,
	}
	ctx.JSON(http.StatusOK,schema.SuccessResponse(responseData,"Access token has been refreshed"))

}

func (a *authHandler) VerifyEmail(ctx *gin.Context){
	var reqCreds *schema.VerifyEmailRequest
	if err:=ctx.BindJSON(&reqCreds);err!=nil{
		resp:=utils.NewAppError(http.StatusBadRequest,"BAD_REQUEST","Request body is invalid, ensure email and verification code are valid",nil)
		_=ctx.Error(resp)
		return
	}
	isVerified,message,_:=a.authService.VerifyEmail(ctx,reqCreds)
	
	if !isVerified{
		resp:=utils.NewAppError(http.StatusUnauthorized,"REQ_UNAUTHORIZED",message,nil)
		_=ctx.Error(resp)
		return
	}

	job_data:=&utils.EmailJob{
		To: reqCreds.Email,
		From: config.AppCfgs.Resend.AppDomain,
		Subj:"Account verification successfully",
		Body: `
		<p>We are pleased to inform you that your email address has been successfully verified.</p>
        <p>You can now access all the features of your account.</p>
        <p>Thank you for joining us!</p>
			<p style="font-size: 12px; color: #888;">Need help? Contact our support team at support@example.com</p>
		`,
	}
	utils.EnqueueEmail(ctx,job_data)
	ctx.JSON(http.StatusOK,schema.SuccessResponse(nil,message))
}

//Access Token Validation godoc
//@Summary Checks if current access token is valid.
//@Descrition This handlers validate wheather current session is active or expired. The handler is protected with JWT middleware. So, middleware will set userID in request context only if given access token is valid.
//@Accept json
//@Product json
//@Security ApiKeyAuth
//@Success 200 {object} schema.SuccessResponseSchema "Token is valid"
//@Failure 401 {object} schema.ErrorResponseSchema "Invalid or missing token"
//@Router /auth/validate [get]
func (a *authHandler) Validate(ctx *gin.Context) {
	userID,exists:=ctx.Get("userID")
	if !exists{
		_=ctx.Error(utils.NewAppError(400,"INVALID_REQUEST","User ID is missing !!",nil))
		return
	}
	
	logger.Log.Info("Attempted Token validation: ",zap.String("userID",userID.(string)),zap.Time("time",time.Now()))

	response:=map[string]string{
		"userID":userID.(string),
	}
	ctx.JSON(http.StatusOK,schema.SuccessResponse(response,"Token is valid"))
}


/*
	Test handlers to test jwt route protection
*/

func (a * authHandler) Profile(ctx *gin.Context){
	userID,exists:=ctx.Get("userID")
	if !exists{
		_=ctx.Error(utils.NewAppError(400,"INVALID_REQUEST","User ID is missing !!",nil))
		return
	}
	logger.Log.Info("Logged user",zap.String("userID",userID.(string)),zap.Time("time",time.Now()))
	response:=map[string]string{
		"userID":userID.(string),
	}
	ctx.JSON(http.StatusOK,schema.SuccessResponse(response,"You are an authorized user"))
}

