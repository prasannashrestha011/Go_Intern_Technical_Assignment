package chimiddlewares

import (
	"main/internal/schema"
	"main/internal/utils"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

func RateLimit(limit rate.Limit,burst int) func(http.Handler)http.Handler{

	return func (next http.Handler)http.Handler  {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ip:=r.RemoteAddr
				allowed,retryAfter:=utils.AllowRequest(ip,limit,burst)
				if !allowed{
					utils.JsonResponseWriter(w,http.StatusTooManyRequests,schema.Response{
						Success: false,
						Error: &schema.ErrorDetail{
							Code: "TOO_MANY_REQUESTS",
							Message: "Too many requests, try again after "+retryAfter.Truncate(time.Second).String(),
						},
					})
					return
				}
				next.ServeHTTP(w,r)
		})
	}

	
}