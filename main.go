package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
)

func main() {
	gin.DefaultWriter = colorable.NewColorableStdout()
	router := gin.Default()
	router.SetTrustedProxies([]string{"192.168.1.2"})
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})
	router.Run() // Listen for "/" on port 8080
}

type Recipe struct {
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}
