package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/config"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Admin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func AdminLogin(c *gin.Context) {
	var req models.RequestAdmin
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin, err := repository.GetAdmin()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateRandomString(32)
	expiresAt := time.Now().Add(1 * time.Hour)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := repository.CreateAdminSession(token, expiresAt); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("admin_session", token, expiresAt.Second(), "/", config.URL.Hostname(), config.URL.Scheme == "https", true)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func AdminPassword(c *gin.Context) {
	var req models.RequestAdminPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin, err := repository.GetAdmin()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.OldPassword)); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := repository.ChangeAdminPassword(string(hashedPassword)); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	repository.DeleteAdminSession()

	c.SetCookie("admin_session", "", -1, "/", config.URL.Hostname(), config.URL.Scheme == "https", true)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func AdminLogout(c *gin.Context) {
	repository.DeleteAdminSession()

	c.SetCookie("admin_session", "", -1, "/", config.URL.Hostname(), config.URL.Scheme == "https", true)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
