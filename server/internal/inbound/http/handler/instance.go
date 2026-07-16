package handler

import (
	"net/http"

	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/dto/request"
	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/dto/response"
	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/utils"
	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InstanceHandler struct {
	instance *usecase.InstanceUseCase
}

func NewInstanceHandler(instance *usecase.InstanceUseCase) InstanceHandler {
	return InstanceHandler{
		instance: instance,
	}
}

func (h *InstanceHandler) CreateForMe(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, "me")
	if err != nil {
		utils.Error(c, err)
		return
	}

	var req request.CreateInstance
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}
	instance, err := req.IntoInstance(me.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}

	createdInstance, err := h.instance.CreateInstance(c.Request.Context(), instance)
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.InstanceToResponse(createdInstance)
	c.JSON(http.StatusCreated, resp)
}

func (h *InstanceHandler) ListForMe(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, "me")
	if err != nil {
		utils.Error(c, err)
		return
	}

	instances, err := h.instance.ListInstancesByFilter(c.Request.Context(), model.InstanceFilter{UserID: &me.ID})
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.InstancesToResponse(instances)
	c.JSON(http.StatusOK, resp)
}

func (h *InstanceHandler) DeleteForMe(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, "me")
	if err != nil {
		utils.Error(c, err)
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}

	if err := h.instance.DeleteInstanceByUserID(c.Request.Context(), id, me.ID); err != nil {
		utils.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Ok)
}
