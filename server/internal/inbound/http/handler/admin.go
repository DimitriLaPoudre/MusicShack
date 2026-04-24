package handler

import (
	"net/http"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/inbound/http/dto/request"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/inbound/http/dto/response"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/inbound/http/utils"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/setup/config"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type AdminHandler struct {
	l        *zerolog.Logger
	cfgHTTP  *config.HTTPConfig
	cfgAdmin *config.AdminConfig
	admin    *usecase.AdminUseCase
}

func NewAdminHandler(l *zerolog.Logger, cfg *config.Config, admin *usecase.AdminUseCase) AdminHandler {
	return AdminHandler{
		l:        l,
		cfgHTTP:  &cfg.HTTP,
		cfgAdmin: &cfg.Admin,
		admin:    admin,
	}
}

func (h *AdminHandler) Login(c *gin.Context) {
	var req request.LoginAdmin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}

	tkn, err := h.admin.Login(c.Request.Context(), req.Password)
	if err != nil {
		utils.Error(c, err, h.l)
		return
	}

	c.SetCookie("admin_token", tkn, int(h.cfgAdmin.TokenDuration.Seconds()), "/api/admin", "", h.cfgHTTP.HTTPS, true)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *AdminHandler) ChangeAdminPassword(c *gin.Context) {

}
