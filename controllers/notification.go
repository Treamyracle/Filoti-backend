package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"filoti-backend/config"
	"filoti-backend/models"

	"github.com/gin-gonic/gin"
)

func GetNotifications(c *gin.Context) {

	var notifications []models.Notification

	if err := config.DB.Preload("Post").Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications: " + err.Error()})
		return
	}

	var notificationsToReturn []gin.H
	for _, notif := range notifications {

		notificationType := "info"
		iconColor := "bg-blue-500"
		if notif.IsRead {
			iconColor = "bg-gray-400"
		}
		if contains(notif.Message, "baru dibuat") {
			notificationType = "new_post"
			iconColor = "bg-green-500"
		} else if contains(notif.Message, "klaim") || contains(notif.Message, "ambil") {
			notificationType = "claim"
			iconColor = "bg-purple-500"
		} else if contains(notif.Message, "update") || contains(notif.Message, "status") {
			notificationType = "update"
			iconColor = "bg-orange-500"
		}

		notificationsToReturn = append(notificationsToReturn, gin.H{
			"id":         notif.ID,
			"post_id":    notif.PostID,
			"message":    notif.Message,
			"is_read":    notif.IsRead,
			"created_at": notif.CreatedAt,
			"time":       formatTimeAgo(notif.CreatedAt),
			"type":       notificationType,
			"iconColor":  iconColor,
			"post_title": notif.Post.Title,
		})
	}

	if notificationsToReturn == nil {
		notificationsToReturn = []gin.H{}
	}

	c.JSON(http.StatusOK, notificationsToReturn)
}

func formatTimeAgo(t time.Time) string {

	diff := time.Since(t)
	if diff < time.Minute {
		return "Baru saja"
	} else if diff < time.Hour {
		return fmt.Sprintf("%d menit lalu", int(diff.Minutes()))
	} else if diff < 24*time.Hour {
		return fmt.Sprintf("%d jam lalu", int(diff.Hours()))
	} else if diff < 30*24*time.Hour {
		return fmt.Sprintf("%d hari lalu", int(diff.Hours()/24))
	}
	return t.Format("02 Jan 2006")
}

func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
