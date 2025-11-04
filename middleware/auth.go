package middleware

import (
	"net/http"

	"filoti-backend/config"
	"filoti-backend/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		uid := session.Get("id")
		if uid == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		userID, ok := uid.(uint)
		if !ok {

			if tmpInt, ok2 := uid.(int); ok2 {
				userID = uint(tmpInt)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
				c.Abort()
				return
			}
		}

		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
