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

func NewRecipe(c *gin.Context) {
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

// example json for recipe
// {
//     "name": "Homemade Pizza",
//     "tags": [
//         "italian",
//         "pizza",
//         "dinner"
//     ],
//     "ingredients": [
//         "1 1/2 cups (355 ml) warm water (105°F-115°F)",
//         "1 package (2 1/4 teaspoons) of active dry yeast",
//         "3 3/4 cups (490 g) bread flour",
//         "feta cheese, firm mozzarella cheese, grated"
//     ],
//     "instructions": [
//         "Step 1.",
//         "Step 2.",
//         "Step 3."
//     ]
// }
