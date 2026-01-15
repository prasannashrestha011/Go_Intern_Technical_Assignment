package ginmiddlewares

import (
	"fmt"
	"main/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func MetricMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		start := time.Now()
		ctx.Next()
		duration := time.Since(start)
		utils.HttpRequestDuration.Observe(duration.Seconds())
		utils.HttpRequestsTotal.WithLabelValues(ctx.Request.Method, ctx.FullPath(), fmt.Sprintf("%d", ctx.Writer.Status())).Inc()
	}
}
