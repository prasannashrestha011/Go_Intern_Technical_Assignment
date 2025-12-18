package ginmiddlewares

import (
	"main/internal/logger"
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
			ctx.JSON(http.StatusForbidden,gin.H{
				"error":"Authorization token not provided",
			})
			ctx.Abort()
			return
		}
		tokenStr:=parts[1]
		token,err:=utils.ValidateJWT(tokenStr)
		if err!=nil{
			logger.Log.Error("Token expiration error: ",zap.Error(err))
			ctx.JSON(http.StatusForbidden,gin.H{
				"error":"Token expired",
			})
			ctx.Abort()
			return
		}
		userID,err:=utils.GenerateUserIDFromToken(token)
		if err!=nil{
			ctx.JSON(http.StatusForbidden,gin.H{
				"error":"Invalid access token",
			})
			ctx.Abort()
			return
		}
		ctx.Set("userID",userID)
		ctx.Next()
	
	}
}