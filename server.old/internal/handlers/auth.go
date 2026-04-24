package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/config"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/services"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var req models.RequestUserLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		err := fmt.Errorf("c.ShouldBindJSON: %w", err)
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := repository.GetUserByUsername(req.Username)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		err := fmt.Errorf("bcrypt.CompareHashAndPassword: %w", err)
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	expiresIn := (24 * time.Hour)

	session, err := services.CreateUserSession(user.ID, expiresIn)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.SetCookie("user_session", session.Token, int(expiresIn.Seconds()), "/", "", config.HTTPS, true)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func Logout(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.Cookie("user_session")
	if err != nil {
		err := fmt.Errorf("c.Cookie: %w", err)
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repository.DeleteUserSession(userId, token); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("user_session", "", -1, "/", "", config.HTTPS, true)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
