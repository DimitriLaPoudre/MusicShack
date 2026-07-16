package handler

import (
	"net/http"

	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/dto/request"
	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/dto/response"
	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/utils"
	"github.com/DimitriLaPoudre/MusicShack/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	user *usecase.UserUseCase
}

func NewUserHandler(user *usecase.UserUseCase) UserHandler {
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
	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}

	user, err := h.user.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.UserToResponse(user)
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
	var req request.UpdateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}
	id := c.Param("id")
	user, err := req.IntoPartialUser(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}

	updatedUser, err := h.user.UpdateUser(c.Request.Context(), user)
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.UserToResponse(updatedUser)
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}

	if err := h.user.DeleteUser(c.Request.Context(), userID); err != nil {
		utils.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Ok)
}
