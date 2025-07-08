package routes

import (
	"gin-blog-api/controllers"
	"gin-blog-api/middlewares"
	

	"github.com/gin-gonic/gin"
)

func User(r *gin.Engine) {
	api := r.Group("/user")
	{
		api.POST("/signup", controllers.Signup)
		api.POST("/login", controllers.Login)

		api.GET("/posts", controllers.GetAllPosts)
		api.GET("/users/:id/posts", controllers.GetPostsByUser)

		api.GET("/posts/:id", controllers.GetPostByID)

		api.GET("/users/:id/follow-counts", controllers.GetUserFollowCounts)
		api.GET("/users/:id/profile", controllers.GetUserProfile)
		api.GET("/users/:id/followers/count", controllers.GetFollowersCount)
		api.GET("/users/:id/following/count", controllers.GetFollowingCount)

		protected := api.Group("/protected")
		protected.Use(middlewares.JWTAuthMiddleware())

		protected.PUT("/users/me", controllers.UpdateUserProfile)
		protected.PUT("/users/me/password", controllers.ChangePassword)

		protected.POST("/posts", controllers.CreatePost)
		protected.PUT("/posts/:id", controllers.UpdatePost)
		protected.DELETE("/posts/:id", controllers.DeletePost)
		protected.GET("/feed", controllers.GetFeed)
		protected.GET("/users/me/liked-posts", controllers.GetLikedPosts)

		protected.POST("/posts/:id/comments", controllers.CreateComment)
		protected.PUT("/comments/:id", controllers.UpdateComment)
		protected.DELETE("/comments/:id", controllers.DeleteComment)

		protected.POST("/posts/:id/like", controllers.ToggleLike)

		protected.GET("/me", controllers.GetCurrentUser)

		protected.POST("/users/:id/follow", controllers.FollowUser)

		protected.GET("/notifications", controllers.GetNotifications)
		protected.PUT("/notifications/:id/read", controllers.MarkNotificationAsRead)

	}

}

// func WebSocketRoutes(r *gin.Engine) {
// 	r.GET("/ws/notifications", ws.NotificationSocket)
// }


func AuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.GET("/google/login", controllers.GoogleLogin)
		auth.GET("/google/callback", controllers.GoogleCallback)
		// Diğer auth route'lar...
	}
}
