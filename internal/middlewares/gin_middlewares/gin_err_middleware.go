package ginmiddlewares

import (
	"main/internal/logger"
	"main/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context) {

		ctx.Next()
		if len(ctx.Errors)==0{
			return
		}

		err:=ctx.Errors.Last().Err

		switch e:=err.(type){
			case *utils.AppError:
			logger.Log.Error("Request failed",
				zap.Int("status",e.Code),
				zap.String("message",e.Message))
			ctx.AbortWithStatusJSON(e.Code,gin.H{
				"error":e.Message,
			})
		default:
			logger.Log.Error("Internal server error",zap.Error(err))
			ctx.JSON(http.StatusInternalServerError,gin.H{
				"error":"Internal server error",
			})

		}

	}
}