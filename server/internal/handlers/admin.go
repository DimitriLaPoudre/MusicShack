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
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Admin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func AdminLogin(c *gin.Context) {
	var req models.RequestAdmin
	if err := c.ShouldBindJSON(&req); err != nil {
		err := fmt.Errorf("c.ShouldBindJSON: %w", err)
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin, err := repository.GetAdmin()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		err := fmt.Errorf("bcrypt.CompareHashAndPassword: %w", err)
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	expiresIn := (1 * time.Hour)
	if token, err := services.CreateAdminSession(expiresIn); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.SetCookie("admin_session", token, int(expiresIn.Seconds()), "/", "", config.HTTPS, true)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func AdminPassword(c *gin.Context) {
	var req models.RequestAdminPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		err := fmt.Errorf("c.ShouldBindJSON: %w", err)
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.ValidatePassword(req.NewPassword); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	admin, err := repository.GetAdmin()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.OldPassword)); err != nil {
		err := fmt.Errorf("bcrypt.CompareHashAndPassword: %w", err)
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		err := fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := repository.ChangeAdminPassword(string(hashedPassword)); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	services.AdminLogout(c)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func AdminLogout(c *gin.Context) {
	services.AdminLogout(c)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
