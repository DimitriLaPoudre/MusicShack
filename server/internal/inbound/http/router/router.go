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
	adminMW gin.HandlerFunc,
	// userMW gin.HandlerFunc,

	// authH *handler.AuthHandler,
	adminH *handler.AdminHandler,
	userH *handler.UserHandler,
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

		usersGroup := api.Group("/users")
		{
			usersGroup.Use(middleware.RateLimiter(time.Minute, 25))
			usersGroup.Use(adminMW)
			usersGroup.POST("/", userH.Create)
			usersGroup.GET("/", userH.List)
			usersGroup.GET("/:id", userH.GetByID)
			usersGroup.PUT("/:id", userH.Update)
			usersGroup.DELETE("/:id", userH.Delete)
		}

		// authGroup := v1.Group("/auth")
		// {
		// 	authGroup.POST("/signup", middleware.RateLimiter(time.Minute, 5), guestMW, authH.SignupLogin)
		// 	authGroup.POST("/login", middleware.RateLimiter(time.Minute, 10), guestMW, authH.Login)
		// 	authGroup.DELETE("/logout", middleware.RateLimiter(time.Minute, 10), authH.Logout)
		// 	authGroup.PUT("/refresh", middleware.RateLimiter(time.Minute, 10), authH.RefreshToken)
		// }

	}
}

func healthz(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
