package handler

import (
	"net/http"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/inbound/http/dto/request"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/inbound/http/dto/response"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type AdminHandler struct {
	l     *zerolog.Logger
	admin *usecase.AdminUseCase
}

func NewAdminHandler(l *zerolog.Logger, admin *usecase.AdminUseCase) AdminHandler {
	return AdminHandler{
		l:     l,
		admin: admin,
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
		c.JSON(http.StatusInternalServerError, response.NewError(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tkn})
}

func (h *AdminHandler) ChangeAdminPassword(c *gin.Context) {

}
