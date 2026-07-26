package handler

import (
	"net/http"

	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/dto/request"
	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/dto/response"
	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/macro"
	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/utils"
	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FollowHandler struct {
	follow *service.FollowService
}

func NewFollowHandler(follow *service.FollowService) FollowHandler {
	return FollowHandler{
		follow: follow,
	}
}

func (h *FollowHandler) Add(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
	if err != nil {
		utils.Error(c, err)
		return
	}

	var req request.CreateFollow
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}
	follow, err := req.IntoFollow(me.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}

	createdFollow, err := h.follow.CreateFollow(c.Request.Context(), follow)
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.FollowToResponse(createdFollow)
	c.JSON(http.StatusOK, resp)
}

func (h *FollowHandler) List(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
	if err != nil {
		utils.Error(c, err)
		return
	}

	follows, err := h.follow.ListFollowsByFilter(c.Request.Context(), model.FollowFilter{UserID: &me.ID})
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.FollowsToResponse(follows)
	c.JSON(http.StatusOK, resp)
}

func (h *FollowHandler) Delete(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
	if err != nil {
		utils.Error(c, err)
		return
	}

	idStr := c.Param(macro.ID)
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}

	if err := h.follow.DeleteFollowByUserID(c.Request.Context(), id, me.ID); err != nil {
		utils.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Ok)
}
