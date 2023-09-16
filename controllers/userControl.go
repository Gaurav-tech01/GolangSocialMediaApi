package controllers

import (
	"context"
	"net/http"
	"social-media/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUsers(c *gin.Context) {
	userID := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(userID)
	var user models.User
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": user})
		return
	}
	c.JSON(http.StatusOK, user)
}

func GetUserFriends(c *gin.Context) {
	userID := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(userID)
	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	var friend []string = user.Friends
	var allFriend []models.User
	for i := 0; i < len(friend); i++ {
		var found models.User
		id, _ := primitive.ObjectIDFromHex(friend[i])
		err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&found)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": friend})
			return
		}
		allFriend = append(allFriend, found)
	}
	c.JSON(http.StatusOK, allFriend)
}

func AddRemoveFriend(c *gin.Context) {
	userID := c.Param("id")
	friendID := c.Param("friendId")
	id, _ := primitive.ObjectIDFromHex(userID)
	fid, _ := primitive.ObjectIDFromHex(friendID)
	var user models.User
	var friend models.User
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Heelo"})
		return
	}
	err = collection.FindOne(context.Background(), bson.M{"_id": fid}).Decode(&friend)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	f := 0
	f1 := 0
	for i, v := range user.Friends {
		if v == friendID {
			copy(user.Friends[i:], user.Friends[i+1:])
			user.Friends = user.Friends[:len(user.Friends)-1]
			f = 1
			break
		}
	}
	if f == 0 {
		user.Friends = append(user.Friends, friendID)
	}
	for i, v := range friend.Friends {
		if v == userID {
			copy(friend.Friends[i:], friend.Friends[i+1:])
			friend.Friends = friend.Friends[:len(friend.Friends)-1]
			f1 = 1
			break
		}
	}
	if f1 == 0 {
		friend.Friends = append(friend.Friends, userID)
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"friends": user.Friends}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User friend not updated"})
		return
	}
	filter = bson.M{"_id": fid}
	update = bson.M{"$set": bson.M{"friends": friend.Friends}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "friend friend not updated"})
		return
	}
	c.JSON(http.StatusOK, result)
}
