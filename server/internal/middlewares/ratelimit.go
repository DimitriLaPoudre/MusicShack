package middlewares

import (
	"github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimiter(formatted string) gin.HandlerFunc {
	rate, _ := limiter.NewRateFromFormatted(formatted)
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	return ginlimiter.NewMiddleware(instance, ginlimiter.WithLimitReachedHandler(
		func(c *gin.Context) {
			c.JSON(429, gin.H{"error": "Too many requests, try later"})
			c.Abort()
		}))
}
