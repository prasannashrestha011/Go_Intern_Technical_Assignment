package ginmiddlewares

import (
	"main/internal/logger"
	"main/internal/schema"
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
				zap.Int("status code",e.StatusCode),
				zap.String("error message",e.Message),
				zap.String("error details",e.Details))
			ctx.AbortWithStatusJSON(e.StatusCode,schema.Response{
				Success: false,
				Message: e.Details,
				Error: &schema.ErrorDetail{
					Message: e.Message,
					Details: e.Details,
					Code: e.Code,
				},
			})
		default:
			logger.Log.Error("Internal server error",zap.Error(err))
			ctx.JSON(http.StatusInternalServerError,schema.Response{
				Success: false,
				Message: "Unknown error occured in the server",
				Error: &schema.ErrorDetail{
					Code: "INTERNAL_SERVER_ERR",
					Message: "Internal server error",
					Details: err.Error(),
				},
			})

		}

	}
}