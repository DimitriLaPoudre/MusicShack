package middlewares

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimiter(formatted string) gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted(formatted)
	if err != nil {
		log.Fatal("limiter.NewRateFromFormatted:", err)
	}
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	return ginlimiter.NewMiddleware(instance, ginlimiter.WithLimitReachedHandler(
		func(c *gin.Context) {
			err := errors.New("too many requests, try later")
			log.Println(err)
			c.JSON(429, gin.H{"error": err})
			c.Abort()
		}))
}
