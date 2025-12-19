package ginmiddlewares

import (
	"main/internal/logger"
	"main/internal/schema"
	"main/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GinJWTMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		authHeader:=ctx.GetHeader("Authorization")
		parts:=strings.Split(authHeader, "Bearer ")

		if len(parts)!=2 || strings.TrimSpace(parts[1])==""{
			ctx.JSON(http.StatusUnauthorized,schema.ErrorResponse("BAD_REQUEST","Invalid access token","Either the authorization token is invalid or missing from authorization header."))
			ctx.Abort()
			return
		}
		tokenStr:=parts[1]
		token,err:=utils.ValidateJWT(tokenStr)
		if err!=nil{
			logger.Log.Error("Token expiration error: ",zap.Error(err))
			ctx.JSON(http.StatusUnauthorized,schema.ErrorResponse("BAD_REQUEST","Access token expired, use refresh token to revise the session life.",""))
			ctx.Abort()
			return
		}
		userID,err:=utils.GenerateUserIDFromToken(token)
		if err!=nil{
			ctx.JSON(http.StatusUnauthorized,schema.ErrorResponse("TOKEN_GEN_FAILURE","Invalid access token, Unable to extract userID from the claims",""))
			ctx.Abort()
			return
		}
		ctx.Set("userID",userID)
		ctx.Next()
	
	}
}