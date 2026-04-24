package router

import (
	"net/http"
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/inbound/http/handler"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/inbound/http/middleware"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/setup/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func New(
	app *gin.Engine,
	l *zerolog.Logger,
	cfg *config.Config,

	// authMW gin.HandlerFunc,
	// guestMW gin.HandlerFunc,
	// adminMW gin.HandlerFunc,
	// userMW gin.HandlerFunc,
	//
	// userH *handler.UserHandler,
	// authH *handler.AuthHandler,
	adminH *handler.AdminHandler,
) {
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))
	gin.Recovery()

	app.GET("/healthz", healthz)

	api := app.Group("/api")
	{
		adminGroup := api.Group("/admin")
		{
			adminGroup.POST("/login", middleware.RateLimiter(time.Minute, 3), adminH.Login)
		}

		// authGroup := v1.Group("/auth")
		// {
		// 	authGroup.POST("/signup", middleware.RateLimiter(time.Minute, 5), guestMW, authH.SignupLogin)
		// 	authGroup.POST("/login", middleware.RateLimiter(time.Minute, 10), guestMW, authH.Login)
		// 	authGroup.DELETE("/logout", middleware.RateLimiter(time.Minute, 10), authH.Logout)
		// 	authGroup.PUT("/refresh", middleware.RateLimiter(time.Minute, 10), authH.RefreshToken)
		// }

		// usersGroup := v1.Group("/users")
		// {
		// 	usersGroup.Use(authMW, adminMW)
		// 	usersGroup.POST("/", middleware.RateLimiter(time.Minute, 25), userH.Create)
		// 	usersGroup.GET("/", middleware.RateLimiter(time.Minute, 100), userH.List)
		// 	usersGroup.GET("/:id", middleware.RateLimiter(time.Minute, 100), userH.GetByID)
		// 	usersGroup.PUT("/:id", middleware.RateLimiter(time.Minute, 25), userH.Update)
		// 	usersGroup.DELETE("/:id", middleware.RateLimiter(time.Minute, 25), userH.Delete)
		// }
	}
}

func healthz(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
