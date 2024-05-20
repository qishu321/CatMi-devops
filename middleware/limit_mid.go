package middleware

import (
	"CatMi-devops/response"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"time"
)

func RateLimitMiddleware(fillInterval time.Duration, capacity int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucket(fillInterval, capacity)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			response.Fail(c, nil, "访问限流")
			c.Abort()
			return
		}
		c.Next()
	}
}
