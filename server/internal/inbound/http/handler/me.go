package handler

import (
	"net/http"

	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/dto/request"
	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/dto/response"
	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/utils"
	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/service"
	"github.com/gin-gonic/gin"
)

type MeHandler struct {
	user *service.UserService
}

func NewMeHandler(user *service.UserService) MeHandler {
	return MeHandler{
		user: user,
	}
}

func (h *MeHandler) Get(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, "me")
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.UserToResponse(me)
	c.JSON(http.StatusOK, resp)
}

func (h *MeHandler) Update(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, "me")
	if err != nil {
		utils.Error(c, err)
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

	if _, err := h.user.UpdateUser(c, partialUser); err != nil {
		utils.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Ok)
}
