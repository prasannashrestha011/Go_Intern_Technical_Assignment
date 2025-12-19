package chimiddlewares

import (
	"context"
	"main/internal/logger"
	"main/internal/utils"
	"net/http"
	"strings"

	"go.uber.org/zap"
)
type userContextKey string
const userKEY=userContextKey("userID")

func JWTAuthMiddleware(next http.Handler)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader:=r.Header.Get("Authorization")
		parts:=strings.Split(authHeader, "Bearer ")
		if len(parts)!=2 || strings.TrimSpace(parts[1])==""{
			appErr:=utils.NewAppError(http.StatusUnauthorized,"REQ_UNAUTHORIZED","Access token missing!!, unable to process the request",nil)
			SetError(w,appErr)
			return
		} 
		tokenStr:=parts[1]
		token,err:=utils.ValidateJWT(tokenStr)
		if err!=nil{
			logger.Log.Info("JWT validation error",zap.Error(err))
			appErr:=utils.NewAppError(http.StatusUnauthorized,"REQ_UNAUTHORIZED","Session expired, please refresh your access token",nil)
			SetError(w,appErr)
			return
		}
		userID,err:=utils.GenerateUserIDFromToken(token)
		if err!=nil{
			logger.Log.Info("JWT ID extraction error",zap.Error(err))
			appErr:=utils.NewAppError(http.StatusUnauthorized,"REQ_UNAUTHORIZED","Invalid access token, failed to retrieve user id from claims",nil)
			SetError(w,appErr)
			return
		}
		ctx:=context.WithValue(r.Context(),userKEY,userID)
		r=r.WithContext(ctx)
		next.ServeHTTP(w,r)
	})
}