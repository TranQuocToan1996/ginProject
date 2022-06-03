package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/TranQuocToan1996/ginProject/model"
	"github.com/TranQuocToan1996/ginProject/recipe"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
)

// init will be executed during the startup of application
func init() {
	recipe.Recipes = make([]*model.Recipe, 0)
	jsonByte, _ := ioutil.ReadFile("recipes.json")
	_ = json.Unmarshal(jsonByte, &recipe.Recipes)
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

	router.POST("/recipes", recipe.AddNewRecipe)
	router.GET("/recipes", recipe.ListRecipes)

	router.Run() // Listen for "/" on port 8080
}
