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

	authMW gin.HandlerFunc,
	adminMW gin.HandlerFunc,
	userMW gin.HandlerFunc,
	guestMW gin.HandlerFunc,

	meH *handler.MeHandler,
	userH *handler.UserHandler,
	authH *handler.AuthHandler,
) {
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))
	gin.Recovery()

	app.GET("/healthz", healthz)

	api := app.Group("/api")
	{
		// adminGroup := api.Group("/admin")
		// {
		// 	adminGroup.POST("/login", middleware.RateLimiter(time.Minute, 5), adminH.Login)
		// 	adminGroup.PUT("/change-password", middleware.RateLimiter(time.Minute, 25), adminH.ChangeAdminPassword)
		// }

		meGroup := api.Group("/me")
		{
			meGroup.Use(middleware.RateLimiter(time.Minute, 100))
			meGroup.Use(authMW)
			meGroup.GET("", meH.Get)
			meGroup.PUT("", meH.Update)

			// instanceGroup := api.Group("/instances")
			// {
			// 	instanceGroup.GET("")
			// }
		}

		usersGroup := api.Group("/users")
		{
			usersGroup.Use(middleware.RateLimiter(time.Minute, 100))
			usersGroup.Use(authMW)
			usersGroup.Use(adminMW)
			usersGroup.POST("", userH.Create)
			usersGroup.GET("", userH.List)
			usersGroup.GET("/:id", userH.GetByID)
			usersGroup.PUT("/:id", userH.Update)
			usersGroup.DELETE("/:id", userH.Delete)
		}

		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", middleware.RateLimiter(time.Minute, 10), guestMW, authH.Login)
			authGroup.DELETE("/logout", middleware.RateLimiter(time.Minute, 10), authMW, authH.Logout)
		}

	}
}

func healthz(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
