package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var apiURL = "https://default.api/"

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		resp, err := http.Get(apiURL)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"context": fmt.Sprintf("GET %s", apiURL), "error": err})
			return
		}
		defer resp.Body.Close()

		var respData struct {
			Version string `json:"HIFI-API"`
			RepoUrl string `json:"Repo"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
			c.JSON(http.StatusOK, gin.H{"message": "API is running...", "api": err})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "API is running...", "api": respData})
		}
	})

	r.Run(":8080")
}
