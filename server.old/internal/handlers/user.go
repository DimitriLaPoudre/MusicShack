package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/config"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(c *gin.Context) {
	var req models.RequestUser
	if err := c.ShouldBindJSON(&req); err != nil {
		err := fmt.Errorf("c.ShouldBindJSON: %w", err)
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.ValidateRequestUser(req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		err := fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := repository.CreateUser(&models.User{Username: req.Username, Password: string(hashPassword)}); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func ListUsers(c *gin.Context) {
	users, err := repository.ListUsers()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		err := fmt.Errorf("strconv.ParseUint: %w", err)
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := repository.GetUserByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := errors.New("user not found")
			log.Println(err)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		err := fmt.Errorf("strconv.ParseUint: %w", err)
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req models.RequestUser
	if err := c.ShouldBindJSON(&req); err != nil {
		err := fmt.Errorf("c.ShouldBindJSON: %w", err)
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oldUser, err := repository.GetUserByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := errors.New("user not found")
			log.Println(err)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = repository.UpdateUser(uint(id), &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := errors.New("user not found")
			log.Println(err)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if req.Username != "" {
		root, err := os.OpenRoot(config.LIBRARY_PATH)
		if err != nil {
			err := fmt.Errorf("os.OpenRoot: %w", err)
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		if err := root.Rename(oldUser.Username, req.Username); err != nil {
			err := fmt.Errorf("root.Rename")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		err := fmt.Errorf("strconv.ParseUint: %w", err)
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = repository.DeleteUser(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := errors.New("user not found")
			log.Println(err)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
