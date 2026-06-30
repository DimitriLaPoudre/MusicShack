package handler

import (
	"net/http"
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/inbound/http/dto/request"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/inbound/http/dto/response"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/inbound/http/utils"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/service"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/setup/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	l          *zerolog.Logger
	cfgHTTP    config.HTTPConfig
	cfgSession config.SessionConfig
	auth       *service.AuthService
}

func NewAuthHandler(l *zerolog.Logger, cfgHTTP config.HTTPConfig, cfgSession config.SessionConfig, auth *service.AuthService) AuthHandler {
	return AuthHandler{
		l:          l,
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

	loginForm := model.LoginForm{
		Username: req.Username,
		Password: req.Password,
		Remember: req.Remember,
	}

	tkn, err := h.auth.Login(c.Request.Context(), &loginForm)
	if err != nil {
		utils.Error(c, err, h.l)
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
		utils.Error(c, err, h.l)
		return
	}

	c.SetCookie(h.cfgSession.CookieName, "", -1, "/api", "", h.cfgHTTP.HTTPS, true)

	if err := h.auth.Logout(c, tkn); err != nil {
		reqID := c.GetString("request_id")
		h.l.Warn().Str("request_id", reqID).Msg(err.Error())
	}

	c.JSON(http.StatusOK, response.Ok)
}
