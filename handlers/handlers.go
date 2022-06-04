package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TranQuocToan1996/ginProject/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecipesHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewRecipesHandler(ctx context.Context, collection *mongo.Collection) *RecipesHandler {
	return &RecipesHandler{
		collection: collection,
		ctx:        ctx,
	}
}

// Recipes holds list of recipes
var Recipes []*models.Recipe

// AddNewRecipe is handler for POST request that include a recipe in JSON
func (handler *RecipesHandler) AddNewRecipe(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err := handler.collection.InsertOne(handler.ctx, recipe)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Error while inserting a new recipe"})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

// ListRecipes returns a list of recipes in JSON format
func (handler *RecipesHandler) ListRecipes(c *gin.Context) {
	cursor, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(handler.ctx)
	recipes := make([]models.Recipe, 0)
	for cursor.Next(handler.ctx) {
		var recipe models.Recipe
		cursor.Decode(&recipe)
		recipes = append(recipes, recipe)
	}
	c.JSON(http.StatusOK, recipes)
}

// UpdateRecipes updates a recipe
func (handler *RecipesHandler) UpdateRecipes(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.UpdateOne(handler.ctx, bson.M{
		"_id": objectId,
	},
		bson.D{{"$set", bson.D{
			{"name", recipe.Name},
			{"instructions", recipe.Instructions},
			{"ingredients", recipe.Ingredients},
			{"tags", recipe.Tags},
		}}})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipe	has been updated"})
}

// DeleteRecipes deletes a recipe
func (handler *RecipesHandler) DeleteRecipes(c *gin.Context) {
	// id := c.Param("id")
	// deleted := false
	// for i := 0; i < len(Recipes); i++ {
	// 	if Recipes[i].ID == id {
	// 		Recipes = append(Recipes[:i], Recipes[i+1:]...)
	// 		deleted = true
	// 		break
	// 	}
	// }
	// if !deleted {
	// 	c.JSON(http.StatusNotFound, gin.H{
	// 		"error": "Recipe not found"})
	// } else {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "Recipe has been deleted",
	// 	})
	// }

}

// SearchRecipes seaches recepi bt the tags
func (handler *RecipesHandler) SearchRecipes(c *gin.Context) {
	// tag := c.Query("tag")
	// listOfRecipes := []*models.Recipe{}
	// for i := 0; i < len(Recipes); i++ {
	// 	for _, t := range Recipes[i].Tags {
	// 		if strings.EqualFold(t, tag) {
	// 			listOfRecipes = append(listOfRecipes, Recipes[i])
	// 		}
	// 	}
	// }
	// c.JSON(http.StatusOK, listOfRecipes)
}
