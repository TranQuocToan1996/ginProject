package recipe

import (
	"net/http"
	"strings"
	"time"

	"github.com/TranQuocToan1996/ginProject/model"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

// Recipes holds list of recipes
var Recipes []*model.Recipe

// AddNewRecipe is handler for POST request that include a recipe in JSON
func AddNewRecipe(c *gin.Context) {
	recipe := &model.Recipe{}
	// Page 47
	if err := c.ShouldBindJSON(recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// recipe.ID = primitive.NewObjectID().Hex()
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()

	Recipes = append(Recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}

// ListRecipes returns a list of recipes in JSON format
func ListRecipes(c *gin.Context) {
	c.JSON(http.StatusOK, Recipes)
}

// UpdateRecipes updates a recipe
func UpdateRecipes(c *gin.Context) {
	id := c.Param("id")
	recipe := &model.Recipe{}
	if err := c.ShouldBindJSON(recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	i := 0
	for ; i < len(Recipes); i++ {
		if Recipes[i].ID == id {
			Recipes[i] = recipe
			break
		}
	}
	if i == len(Recipes)-1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "id not found",
		})
	}
	c.JSON(http.StatusOK, recipe)
}

// DeleteRecipes deletes a recipe
func DeleteRecipes(c *gin.Context) {
	id := c.Param("id")
	delete := false
	for i := 0; i < len(Recipes); i++ {
		if Recipes[i].ID == id {
			Recipes = append(Recipes[:i], Recipes[i+1:]...)
			delete = true
			break
		}
	}
	if !delete {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found"})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Recipe has been deleted",
		})
	}

}

func SearchRecipes(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := []*model.Recipe{}
	for i := 0; i < len(Recipes); i++ {
		found := false
		for _, t := range Recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}
		if found {
			listOfRecipes = append(listOfRecipes, Recipes[i])
		}
	}
	c.JSON(http.StatusOK, listOfRecipes)

}
