package handlers

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/config"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Signup(c *gin.Context) {
	var req models.Signup
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.Username) < 3 || len(req.Username) > 32 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username must be between 3 and 32 characters long"})
		return
	}
	regexUsername := regexp.MustCompile(`^[a-zA-Z0-9_\-.]+$`)
	if !regexUsername.MatchString(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username must only contain alphanumeric characters, _, ., -"})
		return
	}

	if len(req.Password) < 8 || len(req.Password) > 128 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be between 8 and 128 characters long"})
		return
	}
	regexPassword := regexp.MustCompile(`^[a-zA-Z0-9!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'",.<>/?\\|]+$`)
	if !regexPassword.MatchString(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password contains invalid character"})
		return
	}

	_, err := repository.GetUserByUsername(req.Username)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username already used"})
		return
	} else {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := repository.CreateUser(&models.User{Username: req.Username, Password: string(hashPassword)}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func Login(c *gin.Context) {
	var req models.Login
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := repository.GetUserByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "username not found"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	token, err := services.GetTokenForID(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", token, 0, "/", config.URL.Hostname(), config.URL.Scheme == "https", true)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", true, true)
	c.JSON(200, gin.H{"status": "ok"})
}
