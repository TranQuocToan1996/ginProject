package recipe

import (
	"net/http"
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
