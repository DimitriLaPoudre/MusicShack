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

type InstanceHandler struct {
	instance *service.InstanceService
}

func NewInstanceHandler(instance *service.InstanceService) InstanceHandler {
	return InstanceHandler{
		instance: instance,
	}
}

func (h *InstanceHandler) Create(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
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

func (h *InstanceHandler) List(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
	if err != nil {
		utils.Error(c, err)
		return
	}

	_, exists := c.GetQuery("refresh")

	var instances []model.Instance
	if exists {
		instances, err = h.instance.RefreshListByUserID(c.Request.Context(), me.ID)
		if err != nil {
			utils.Error(c, err)
			return
		}
	} else {
		instances, err = h.instance.ListInstancesByFilter(c.Request.Context(), model.InstanceFilter{UserID: &me.ID})
		if err != nil {
			utils.Error(c, err)
			return
		}
	}

	resp := response.InstancesToResponse(instances)
	c.JSON(http.StatusOK, resp)
}

func (h *InstanceHandler) Delete(c *gin.Context) {
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

	if err := h.instance.DeleteInstanceByUserID(c.Request.Context(), id, me.ID); err != nil {
		utils.Error(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
