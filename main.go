package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TranQuocToan1996/ginProject/gcloud"
	"github.com/TranQuocToan1996/ginProject/handlers"
	"github.com/TranQuocToan1996/ginProject/models"
	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var recipesHandler *handlers.RecipesHandler
var authHandler *handlers.AuthHandler

// init will be executed during the startup of application
func init() {
	env := os.Getenv("ENV")
	if len(env) == 0 {
		env = models.Dev + `.env`
	}

	err := godotenv.Load(env)
	if err != nil {
		// Handle case wrong name .env
		err = godotenv.Load(`.env`)
		if err != nil {
			log.Fatal("Fail to load .env")
		}
	}

}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := gcloud.InitFbApp(ctx); err != nil {
		log.Fatal(err)
	}
	_, err := gcloud.GetConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer gcloud.Close()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	// use Ping to check whether the mongo run or not
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to mongodb!")

	collectionRecipes := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	collectionUsers := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")

	// Connect redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost:%v", os.Getenv("REDIS_PORT")),
		Password: "",
		DB:       0,
	})
	redisStatus := redisClient.Ping()
	log.Println(redisStatus)

	// Handler
	recipesHandler = handlers.NewRecipesHandler(ctx, collectionRecipes, redisClient)
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers, redisClient)

	router := gin.Default()
	router.SetTrustedProxies(nil)
	// Load static file in "./static" path and convert into route /static
	router.Static("/static", "./static")

	// Load all templates to memory for faster serve
	router.LoadHTMLGlob("templates/*.html")

	store, _ := redisStore.NewStore(10, "tcp", fmt.Sprintf("localhost:%v", os.Getenv("REDIS_PORT")), "", []byte("secret"))
	router.Use(sessions.Sessions("recipes_api", store))

	router.GET("/", recipesHandler.IndexHandlerHTML)
	router.GET("/recipes/search", recipesHandler.SearchRecipes)
	router.GET("/recipes/search/:id", recipesHandler.SearchRecipeById)
	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/signup", authHandler.RegisterAccount)
	router.POST("/refresh", authHandler.RefreshToken)
	router.POST("/signout", authHandler.SignOut)

	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleware_session())
	{
		authorized.POST("/recipes", recipesHandler.AddNewRecipe)
		authorized.GET("/recipes", recipesHandler.ListRecipes)
		authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipes)
		authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipes)
	}

	// openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout certs/localhost.key -out certs/localhost.crt
	router.RunTLS(":443", "certs/localhost.crt", "certs/localhost.key")
	// router.Run() // Default Listen for "/" on port 8080
}

// Notes
/* docker run -d --name mongodbgin -v mongoGinProject:/data/db -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=password -p 27017:27017 mongo:4.4.3 */

/* docker run -d --name redisinsight --link redisName -p 8002:800ngo)1 redislabs/redisinsight */

// docker run -d --name redisForGin -p 6380:6379 redis:latest

// ab -n 2000 -c 100 -g without-cache.data http://localhost:8080/recipes
