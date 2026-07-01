package handler

import (
	"net/http"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/inbound/http/dto/request"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/inbound/http/dto/response"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/inbound/http/utils"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type InstanceHandler struct {
	l        *zerolog.Logger
	instance *usecase.InstanceUseCase
}

func NewInstanceHandler(l *zerolog.Logger, instance *usecase.InstanceUseCase) InstanceHandler {
	return InstanceHandler{
		l:        l,
		instance: instance,
	}
}

func (h *InstanceHandler) CreateForMe(c *gin.Context) {
	me, err := utils.GetFromContext[*model.User](c, "me")
	if err != nil {
		utils.Error(c, err, h.l)
		return
	}

	var req request.CreateInstance
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}
	instance, err := req.IntoInstance()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}

	createdInstance, err := h.instance.CreateInstance(c.Request.Context(), &instance, me.ID)
	if err != nil {
		utils.Error(c, err, h.l)
		return
	}

	resp := response.InstanceToResponse(createdInstance)
	c.JSON(http.StatusCreated, resp)
}

func (h *InstanceHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}

	instance, err := h.instance.GetInstanceByID(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err, h.l)
		return
	}

	resp := response.InstanceToResponse(instance)
	c.JSON(http.StatusOK, resp)
}

func (h *InstanceHandler) ListForMe(c *gin.Context) {
	me, err := utils.GetFromContext[*model.User](c, "me")
	if err != nil {
		utils.Error(c, err, h.l)
		return
	}

	instances, err := h.instance.ListInstancesByUserID(c.Request.Context(), me.ID)
	if err != nil {
		utils.Error(c, err, h.l)
		return
	}

	resp := response.InstancesToResponse(instances)
	c.JSON(http.StatusOK, resp)
}

func (h *InstanceHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}

	if err := h.instance.DeleteInstance(c.Request.Context(), id); err != nil {
		utils.Error(c, err, h.l)
		return
	}

	c.JSON(http.StatusOK, response.Ok)
}

