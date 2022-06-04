package main

import (
	"context"
	"log"
	"os"

	"github.com/TranQuocToan1996/ginProject/handlers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var recipesHandler *handlers.RecipesHandler

// init will be executed during the startup of application
func init() {
	// Connect to mongodb
	// MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" MONGO_DATABASE=demo go run *.go
	/* 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	   	defer cancel() */
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	// use Ping to check whether the mongo run or not
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to mongodb!")

	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	recipesHandler = handlers.NewRecipesHandler(ctx, collection)

	// result, err := collection.InsertMany(ctx, recipesInterface)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Printf("Inserted %v recipes", len(result.InsertedIDs))
}

func main() {
	/* 	// ctx := context.Background()
	   	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	   	defer cancel()
	   	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	   	if err != nil {
	   		log.Fatal(err)
	   	}
	   	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
	   		log.Fatal(err)
	   	}
	   	log.Println("Connected to mongodb!") */

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	router.POST("/recipes", recipesHandler.AddNewRecipe)
	router.GET("/recipes", recipesHandler.ListRecipes)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipes)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipes)
	router.GET("/recipes/search", recipesHandler.SearchRecipes)

	router.Run() // Default Listen for "/" on port 8080
}

/* docker run -d --name mongodbgin -v mongoGinProject:/data/db -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=password -p 27017:27017 mongo:4.4.3 */
