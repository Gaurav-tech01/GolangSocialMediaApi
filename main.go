package main

import (
	"os"
	"social-media/controllers"
	"social-media/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	//User Routes
	router.POST("/signup", controllers.Signup)
	router.POST("/login", controllers.Login)
	router.GET("/user/:id", middleware.RequireAuth, controllers.GetUsers)
	router.GET("/user/:id/friends", middleware.RequireAuth, controllers.GetUserFriends)
	router.PATCH("/user/:id/friend/:friendId", middleware.RequireAuth, controllers.AddRemoveFriend)

	//Post Routes
	router.POST("/posts", middleware.RequireAuth, controllers.CreatePost)
	router.GET("/getPost", middleware.RequireAuth, controllers.GetFeedPosts)
	router.GET("/:userId/posts", middleware.RequireAuth, controllers.GetUserPosts)
	router.PATCH("/:id/like", middleware.RequireAuth, controllers.LikePost)

	router.Run(":" + port)
}
