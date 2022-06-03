package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/TranQuocToan1996/ginProject/model"
	"github.com/TranQuocToan1996/ginProject/recipe"
	"github.com/gin-gonic/gin"
)

// init will be executed during the startup of application
func init() {
	recipe.Recipes = []*model.Recipe{}
	jsonByte, _ := ioutil.ReadFile("recipes.json")
	_ = json.Unmarshal(jsonByte, &recipe.Recipes)
}

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	router.POST("/recipes", recipe.AddNewRecipe)
	router.GET("/recipes", recipe.ListRecipes)
	router.PUT("/recipes:id", recipe.UpdateRecipes)
	router.DELETE("/recipes:id", recipe.DeleteRecipes)
	router.GET("/recipes/search", recipe.SearchRecipes)

	router.Run() // Listen for "/" on port 8080
}
