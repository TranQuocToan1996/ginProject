package main

import (
	"github.com/TranQuocToan1996/ginProject/model"
	"github.com/TranQuocToan1996/ginProject/recipe"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
)

// init will be executed during the startup of application
func init() {
	recipe.Recipes = make([]*model.Recipe, 0)
}

func main() {
	gin.DefaultWriter = colorable.NewColorableStdout()
	router := gin.Default()
	router.SetTrustedProxies([]string{"192.168.1.2"})

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	router.POST("/recipes", recipe.NewRecipe)

	router.Run() // Listen for "/" on port 8080
}
