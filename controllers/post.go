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

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var postCollection *mongo.Collection = database.DBinit("post")

func CreatePost(c *gin.Context) {
	err := godotenv.Load(".env")
	//Get the email/pass off req body
	file, err := c.FormFile("post_picture")
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
		UserId      string
		Description string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	var user models.User
	id, _ := primitive.ObjectIDFromHex(body.UserId)
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	post := models.Post{UserId: body.UserId, First_name: user.First_name, Last_name: user.Last_name, Location: user.Location, Description: body.Description, User_picture: user.Profile_picture, Picture_path: baseUrl + "/" + newFileName, Likes: make(map[string]bool), Comments: []string{}}
	_, err = postCollection.InsertOne(context.Background(), post)
	if err != nil {
		msg := fmt.Sprintf("Post was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	var posts []models.Post
	cursor, err := postCollection.Find(context.TODO(), bson.D{})
	defer cursor.Close(context.TODO())
	if cursor != nil {
		for cursor.Next(context.TODO()) {
			var post models.Post
			if err := cursor.Decode(&post); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
			posts = append(posts, post)
		}
		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, posts)
	} else {

		c.JSON(http.StatusOK, gin.H{"post": "no post"})
		return
	}
}

func GetFeedPosts(c *gin.Context) {
	var posts []models.Post
	cursor, err := postCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		msg := fmt.Sprintf("Post was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		posts = append(posts, post)
	}
	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func GetUserPosts(c *gin.Context) {
	userID := c.Param("userId")
	var posts []models.Post
	cursor, err := postCollection.Find(context.TODO(), bson.M{"userid": userID})
	if err != nil {
		msg := fmt.Sprintf("Post was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		posts = append(posts, post)
	}
	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func LikePost(c *gin.Context) {
	postID := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(postID)
	var body struct {
		UserId string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	userId := body.UserId
	var posts models.Post
	err := postCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&posts)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	likes := posts.Likes

	if liked, ok := likes[userId]; ok {
		if liked {
			delete(likes, userId)
		} else {
			likes[userId] = true
		}
	} else {
		likes[userId] = true
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"likes": likes}}
	result, err := postCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User friend not updated"})
		return
	}
	c.JSON(http.StatusOK, result)
}
