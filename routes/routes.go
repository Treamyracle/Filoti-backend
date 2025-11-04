package routes

import (
	"net/http"
	"os"
	"time"

	"filoti-backend/controllers"
	"filoti-backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "secret"
	}
	store := cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
	r.Use(sessions.Sessions("gin_session", store))

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5500",
			"http://127.0.0.1:5500",
			"http://localhost:3000",
			"https://filoti-frontend.vercel.app",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.POST("/guest-login", controllers.GuestLogin)
	r.GET("/locations", controllers.GetUniqueLocations)
	r.GET("/posts", controllers.GetPosts)
	r.GET("/posts/:id", controllers.GetPostByID)

	authorized := r.Group("/")
	authorized.Use(middleware.AuthRequired())
	{

		authorized.GET("/me", controllers.GetCurrentUser)
		authorized.POST("/logout", controllers.Logout)
		authorized.GET("/notifications", controllers.GetNotifications)

		posts := authorized.Group("/posts")
		{
			posts.POST("", controllers.CreatePost)
			posts.PUT("/:id", controllers.UpdatePost)
			posts.DELETE("/:id", controllers.DeletePost)
			posts.PUT("/:id/done", controllers.MarkPostAsDone)
		}
	}

	return r
}
