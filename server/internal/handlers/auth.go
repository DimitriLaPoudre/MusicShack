package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/config"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/services"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func validateUsername(username string) error {
	if len(username) < 3 || len(username) > 32 {
		return fmt.Errorf("validateUsername: %s", "username must be between 3 and 32 characters long")
	}
	_, err := repository.GetUserByUsername(username)
	if err == nil {
		return fmt.Errorf("validateUsername: %s", "username already used")
	} else {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("validateUsername: %w", err)
		}
	}
	return nil
}

func validatePassword(password string) error {
	regexUsername := regexp.MustCompile(`^[a-zA-Z0-9_\-.]+$`)
	if !regexUsername.MatchString(password) {
		return fmt.Errorf("validatePassword: %s", "password must only contain alphanumeric characters, _, ., -")
	}

	if len(password) < 8 || len(password) > 128 {
		return fmt.Errorf("validatePassword: %s", "password must be between 8 and 128 characters long")
	}
	regexPassword := regexp.MustCompile(`^[a-zA-Z0-9!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'",.<>/?\\|]+$`)
	if !regexPassword.MatchString(password) {
		return fmt.Errorf("validatePassword: %s", "password contains invalid character")
	}
	return nil
}

func validateRequestUser(req models.UserRequest) error {
	if err := validateUsername(req.Username); err != nil {
		return fmt.Errorf("validateRequestUser: %w", err)
	}

	if err := validatePassword(req.Password); err != nil {
		return fmt.Errorf("validateRequestUser: %w", err)
	}

	return nil
}

func Signup(c *gin.Context) {
	var req models.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validateRequestUser(req); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := repository.CreateUser(&models.User{Username: req.Username, Password: string(hashPassword)}); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func Login(c *gin.Context) {
	var req models.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := repository.GetUserByUsername(req.Username)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	session, err := services.CreateUserSession(user.ID)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.SetCookie("user_session", session.Token, session.ExpiresAt.Second(), "/", config.URL.Hostname(), config.URL.Scheme == "https", true)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func Logout(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.Cookie("user_session")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repository.DeleteUserSession(userId, token); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("user_session", "", -1, "/", config.URL.Hostname(), config.URL.Scheme == "https", true)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
