package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Info(c *gin.Context) {
	// resp, err := http.Get(ApiURL)
	// if err != nil {
	// 	c.JSON(http.StatusServiceUnavailable, gin.H{"context": "GET " + ApiURL, "error": err})
	// 	return
	// }
	// defer resp.Body.Close()
	//
	// var respData struct {
	// 	Version string `json:"HIFI-API"`
	// 	RepoUrl string `json:"Repo"`
	// }
	// if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
	// 	c.JSON(http.StatusOK, gin.H{"message": "API is running...", "api": err})
	// } else {
	// 	c.JSON(http.StatusOK, gin.H{"message": "API is running...", "api": respData})
	// }

	c.JSON(http.StatusOK, gin.H{"message": "API is running..."})

}
