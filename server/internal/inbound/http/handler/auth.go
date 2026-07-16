package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/dto/request"
	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/dto/response"
	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/utils"
	"github.com/DimitriLaPoudre/MusicShack/internal/setup/config"
	"github.com/DimitriLaPoudre/MusicShack/internal/usecase"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	cfgHTTP    config.HTTPConfig
	cfgSession config.SessionConfig
	auth       *usecase.AuthUseCase
}

func NewAuthHandler(cfgHTTP config.HTTPConfig, cfgSession config.SessionConfig, auth *usecase.AuthUseCase) AuthHandler {
	return AuthHandler{
		cfgHTTP:    cfgHTTP,
		cfgSession: cfgSession,
		auth:       auth,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginForm
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}
	loginForm := req.IntoLoginForm()

	tkn, err := h.auth.Login(c.Request.Context(), loginForm)
	if err != nil {
		utils.Error(c, err)
		return
	}

	var expiresAt time.Duration
	if loginForm.Remember {
		expiresAt = h.cfgSession.RememberExp
	} else {
		expiresAt = h.cfgSession.Exp
	}

	c.SetCookie(h.cfgSession.CookieName, tkn, int(expiresAt.Seconds()), "/api", "", h.cfgHTTP.HTTPS, true)
	c.JSON(http.StatusOK, response.Ok)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	tkn, err := c.Cookie(h.cfgSession.CookieName)
	if err != nil {
		utils.Error(c, err)
		return
	}

	c.SetCookie(h.cfgSession.CookieName, "", -1, "/api", "", h.cfgHTTP.HTTPS, true)

	if err := h.auth.Logout(c, tkn); err != nil {
		requestID, err := utils.GetFromContext[string](c, "request_id")
		if err != nil {
			requestID = "unknown"
		}
		slog.Error("logout failed", slog.String("request_id", requestID), slog.String("err", err.Error()))
	}

	c.JSON(http.StatusOK, response.Ok)
}
