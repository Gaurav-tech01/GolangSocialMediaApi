package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"social-media/database"
	"social-media/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func tokenFunc(email string, c *gin.Context) (bs string, err error) {
	err = godotenv.Load(".env")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": email,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	return tokenString, err
}

var collection *mongo.Collection = database.DBinit("user")

func Signup(c *gin.Context) {
	err := godotenv.Load(".env")
	//Get the email/pass off req body
	file, err := c.FormFile("profile_picture")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ext := filepath.Ext(file.Filename)
	if strings.ToLower(ext) != ".jpg" && strings.ToLower(ext) != ".jpeg" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPEG or JPG files are allowed."})
		return
	}

	// Create a new file on the server to store the uploaded image
	newFileName := filepath.Join("./uploads", file.Filename)
	err = c.SaveUploadedFile(file, newFileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	baseUrl := os.Getenv("BASE_URL")
	var body struct {
		First_name string
		Last_name  string
		Email      string
		Password   string
		Friends    []string
		Location   string
		Occupation string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	//Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
		return
	}
	user := models.User{First_name: body.First_name, Last_name: body.Last_name, Email: body.Email, Password: string(hash), Profile_picture: baseUrl + "/" + newFileName, Friends: body.Friends, Location: body.Location, Occupation: body.Occupation}

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		msg := fmt.Sprintf("User item was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, result)
}

func Login(c *gin.Context) {
	err := godotenv.Load(".env")
	//Get the email and pass off req body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	//Look up requested user
	var found models.User
	err = collection.FindOne(context.Background(), bson.M{"email": body.Email}).Decode(&found)
	if found.Email == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
	}
	//Comapare sent in pass with saved user pass hash
	err = bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}
	token, err := tokenFunc(body.Email, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token not Found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
