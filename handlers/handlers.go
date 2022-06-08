package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/TranQuocToan1996/ginProject/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const noExpirationTimeRedis = 0

type RecipesHandler struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

func NewRecipesHandler(ctx context.Context, collection *mongo.Collection, redisClient *redis.Client) *RecipesHandler {
	return &RecipesHandler{
		collection:  collection,
		ctx:         ctx,
		redisClient: redisClient,
	}
}

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

	// The data cached in memory, so if we install/update new recipe in mongo. The cached not update yet.
	// There are 2 solutions for this situation.
	// First, Set Time TO Live (TTL) for the recipes.
	// Second is delete recipes after install/update new recipe. The redis will load again when we call ListRecipes.
	// In the scope of project, the data not so much. So we choose 2nd solution
	handler.redisClient.Del("recipes")
	log.Println("Removed redis recipes!")

	c.JSON(http.StatusOK, recipe)
}

// ListRecipes returns a list of recipes in JSON format
func (handler *RecipesHandler) ListRecipes(c *gin.Context) {
	recipes := make([]models.Recipe, 0)

	redisVal, err := handler.redisClient.Get("recipes").Result()
	if err == redis.Nil {
		log.Println("Redis nil, Need to query data from mongo!")
		// Query mongo
		cursor, err := handler.collection.Find(handler.ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(handler.ctx)

		for cursor.Next(handler.ctx) {
			var recipe models.Recipe
			cursor.Decode(&recipe)
			recipes = append(recipes, recipe)
		}
		// Redis value has to be a string -> need encode
		data, _ := json.Marshal(recipes)

		handler.redisClient.Set("recipes", string(data), noExpirationTimeRedis)
		c.JSON(http.StatusOK, recipes)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		log.Println("Redis has data, starting query to Redis!")
		json.Unmarshal([]byte(redisVal), &recipes)
		c.JSON(http.StatusOK, recipes)
	}

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

	handler.redisClient.Del("recipes")
	log.Println("Removed redis recipes!")
	c.JSON(http.StatusOK, gin.H{"message": "Recipe	has been updated"})
}

// DeleteRecipes deletes a recipe
func (handler *RecipesHandler) DeleteRecipes(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.DeleteOne(handler.ctx, bson.M{
		"_id": objectId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been deleted"})

}

// SearchRecipeById return 1 recipe by mongo _id
func (handler *RecipesHandler) SearchRecipeById(c *gin.Context) {
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	singleResulst := handler.collection.FindOne(handler.ctx, bson.M{
		"_id": objectId,
	})
	returnRecipe := models.Recipe{}
	err = singleResulst.Decode(&returnRecipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, returnRecipe)
}

// SearchRecipes seaches recipes by the tags
func (handler *RecipesHandler) SearchRecipes(c *gin.Context) {

	recipes := make([]models.Recipe, 0)

	redisVal, err := handler.redisClient.Get("recipes").Result()
	if err == redis.Nil {
		log.Println("Redis nil, Need to query data from mongo!")
		// Query mongo
		cursor, err := handler.collection.Find(handler.ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(handler.ctx)

		for cursor.Next(handler.ctx) {
			var recipe models.Recipe
			cursor.Decode(&recipe)
			recipes = append(recipes, recipe)
		}
		// Redis value has to be a string -> need encode
		data, _ := json.Marshal(recipes)

		handler.redisClient.Set("recipes", string(data), noExpirationTimeRedis)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		log.Println("Redis has data, starting query to Redis!")
		json.Unmarshal([]byte(redisVal), &recipes)
	}

	tag := c.Query("tag")
	if len(tag) == 0 {
		c.JSON(http.StatusOK, recipes)
		return
	}

	listOfRecipes := make([]models.Recipe, 0)
	for i := 0; i < len(recipes); i++ {
		found := false
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}
		if found {
			listOfRecipes = append(listOfRecipes, recipes[i])
		}
	}

	c.JSON(http.StatusOK, listOfRecipes)
}
