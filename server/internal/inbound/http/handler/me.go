package handler

import (
	"net/http"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/inbound/http/dto/request"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/inbound/http/dto/response"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/inbound/http/utils"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type MeHandler struct {
	l    *zerolog.Logger
	user *service.UserService
}

func NewMeHandler(l *zerolog.Logger, user *service.UserService) MeHandler {
	return MeHandler{
		l:    l,
		user: user,
	}
}

func (h *MeHandler) Get(c *gin.Context) {
	me, err := utils.GetFromContext[*model.User](c, "me")
	if err != nil {
		utils.Error(c, err, h.l)
		return
	}

	resp := response.UserToResponse(me)
	c.JSON(http.StatusOK, resp)
}

func (h *MeHandler) Update(c *gin.Context) {
	me, err := utils.GetFromContext[*model.User](c, "me")
	if err != nil {
		utils.Error(c, err, h.l)
		return
	}

	var req request.UpdateMe
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}
	partialUser, err := req.IntoPartialUser(me.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}

	if _, err := h.user.UpdateUser(c, &partialUser); err != nil {
		utils.Error(c, err, h.l)
		return
	}

	c.JSON(http.StatusOK, response.Ok)
}
