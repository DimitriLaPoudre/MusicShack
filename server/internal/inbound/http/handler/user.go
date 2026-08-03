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

type UserHandler struct {
	user *service.UserService
}

func NewUserHandler(user *service.UserService) UserHandler {
	return UserHandler{
		user: user,
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req request.CreateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}
	user, err := req.IntoUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}

	createdUser, err := h.user.CreateUser(c.Request.Context(), user)
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.UserToResponse(createdUser)
	c.JSON(http.StatusCreated, resp)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	targetUser, err := utils.GetFromContext[model.User](c, "target_user")
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.UserToResponse(targetUser)
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.user.ListAllUsers(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.UsersToResponse(users)
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Update(c *gin.Context) {
	targetUser, err := utils.GetFromContext[model.User](c, "target_user")
	if err != nil {
		utils.Error(c, err)
		return
	}

	var req request.UpdateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}
	partialUser := req.IntoPartialUser(targetUser.ID)

	updatedUser, err := h.user.UpdateUser(c.Request.Context(), partialUser)
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.UserToResponse(updatedUser)
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Delete(c *gin.Context) {
	targetUser, err := utils.GetFromContext[model.User](c, "target_user")
	if err != nil {
		utils.Error(c, err)
		return
	}

	if err := h.user.DeleteUser(c.Request.Context(), targetUser.ID); err != nil {
		utils.Error(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
