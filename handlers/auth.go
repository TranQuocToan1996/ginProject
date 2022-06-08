package handlers

import (
	"context"
	"crypto/sha256"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/TranQuocToan1996/ginProject/models"
	"github.com/auth0-community/go-auth0"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/square/go-jose.v2"
)

type AuthHandler struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

func NewAuthHandler(ctx context.Context, collection *mongo.Collection, redisClient *redis.Client) *AuthHandler {
	return &AuthHandler{
		collection:  collection,
		ctx:         ctx,
		redisClient: redisClient,
	}
}

type Claims struct {
	UserName string `json:"username"`
	jwt.StandardClaims
}

type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

func (handler *AuthHandler) SignInHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	h256 := sha256.New()
	cursor := handler.collection.FindOne(handler.ctx, bson.M{
		"username": user.UserName,
		"password": string(h256.Sum([]byte(user.Password))),
	})

	if cursor.Err() != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}
	// JWT
	expirationTime := time.Now().Add(time.Hour)
	claims := &Claims{
		UserName: user.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	jwtOutPut := JWTOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}
	// JWT

	// Session
	// TO use this, need to use AuthMiddleware_session()
	// sessionToken := xid.New().String()
	// session := sessions.Default(c)
	// session.Set("username", user.Username)
	// session.Set("token", sessionToken)
	// session.Save()
	// Session

	c.JSON(http.StatusOK, jwtOutPut)
}

func (handler *AuthHandler) RegisterAccountHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	cursor := handler.collection.FindOne(handler.ctx, bson.M{
		"username": user.UserName,
	})

	if cursor.Err() != nil {
		if cursor.Err() != mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": cursor.Err(),
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errors.New("exist user name"),
			})
		}
		return
	}

	h := sha256.New()
	result, err := handler.collection.InsertOne(handler.ctx, bson.M{
		"username": user.UserName,
		"password": string(h.Sum([]byte(user.Password))),
	})

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":    "register OK",
		"insertedID": result.InsertedID,
	})
}

// RefreshTokenHandler refresh the token
func (handler *AuthHandler) RefreshTokenHandler(c *gin.Context) {
	tokenVal := c.GetHeader("Authorization")
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenVal, claims,
		func(tkn *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	if token == nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
		})
		return
	}
	if time.Until(time.Unix(claims.ExpiresAt, 0)) >= time.Minute {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token is not expired yet",
		})
		return
	}

	expirationTime := time.Now().Add(time.Hour)

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	jwtOutPut := JWTOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}
	c.JSON(http.StatusOK, jwtOutPut)

}

// AuthMiddleware_APIKEY apply APIKEY
func (handler *AuthHandler) AuthMiddleware_APIKEY() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-API-KEY") != os.Getenv("X_API_KEY") {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()
	}
}

// func (handler *AuthHandler) SignOutHandler(c*gin.Context) {
// 	session := sessions.Default(c)
// 	session.Clear()
// 	session.Save()
// 	c.JSON(http.StatusOK, gin.H{"message":"Signed out..."})
// }

// AuthMiddleware_session obtain the token from the request cookie. If
// the cookie is not set, we return a 403 code (Forbidden) by returning an http.
// StatusForbidden response, as illustrated in the following code snippet
func (handler *AuthHandler) AuthMiddleware_session() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionToken := session.Get("token")
		if sessionToken == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Not logged",
			})
			c.Abort()
		}
		c.Next()
	}
}

// AuthMiddleware JWT
/*  Next, we update the authentication middleware in handler/auth.go to check for the Authorization header instead of the X-API-KEY attribute. The header is then passed to the ParseWithClaims method. It generates a signature using the header and payload from the Authorization header and the secret key. Ten, it verifies if the signature matches the one on the JWT. If not, the JWT is not considered valid, and a 401 status code is returned. The Go implementation is shown here */
func (handler *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenVal := c.GetHeader("Authorization")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenVal, claims,
			func(tkn *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		if token == nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()
	}
}

func (handler *AuthHandler) AuthMiddleware_Auth0() gin.HandlerFunc {
	return func(c *gin.Context) {
		var auth0Domain = "https://" + os.Getenv("AUTH0_DOMAIN") + "/"
		client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: auth0Domain + ".well-known/jwks.json"}, nil)
		configuration := auth0.NewConfiguration(client, []string{os.Getenv("AUTH0_API_IDENTIFIER")}, auth0Domain, jose.RS256)
		validator := auth0.NewValidator(configuration, nil)
		_, err := validator.ValidateRequest(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}
		c.Next()
	}
}
