package ginmiddlewares

import (
	"main/internal/logger"
	"main/internal/schema"
	"main/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

func RateLimit(limit rate.Limit,burst int)gin.HandlerFunc{
	return func(ctx *gin.Context) {
		ip:=ctx.ClientIP()
		
		allowed,retryAfter:=utils.AllowRequest(ip,limit,burst)
		logger.Log.Info("Request throttler",zap.Bool("Is Allowed",allowed))
		if !allowed{

			if retryAfter>0{
				ctx.AbortWithStatusJSON(http.StatusTooManyRequests,schema.Response{
					Success: false,
					Error: &schema.ErrorDetail{
						Code: "TOO_MANY_REQUESTS",
						Message: "Too many requests, try again after "+retryAfter.Truncate(time.Second).String(),
					},
				})
			}
			return
		}
		ctx.Next()
	}
}