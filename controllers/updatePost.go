package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"filoti-backend/config"
	"filoti-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UpdatePostInput struct {
	Title      string `json:"title"`
	Keterangan string `json:"keterangan"`
	Ruangan    string `json:"ruangan"`
	ImageURL   string `json:"image_url"`
}

func UpdatePost(c *gin.Context) {

	uidVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var currentUserID uint
	switch v := uidVal.(type) {
	case uint:
		currentUserID = v
	case int:
		currentUserID = uint(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, currentUserID).Error; err != nil || !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Only administrators can update posts"})
		return
	}

	id := c.Param("id")
	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve post: " + err.Error()})
		return
	}

	var input UpdatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&post).Updates(models.Post{
		Title:      input.Title,
		Keterangan: input.Keterangan,
		Ruangan:    input.Ruangan,
		ImageURL:   input.ImageURL,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully", "post": post})
}

func DeletePost(c *gin.Context) {

	uidVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var currentUserID uint
	switch v := uidVal.(type) {
	case uint:
		currentUserID = v
	case int:
		currentUserID = uint(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, currentUserID).Error; err != nil || !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Only administrators can delete posts"})
		return
	}

	id := c.Param("id")
	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve post: " + err.Error()})
		return
	}

	if err := config.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

type MarkDoneInput struct {
	ClaimerName string `json:"claimer_name" binding:"required"`
	ProofImage  string `json:"proof_image"`
}

func MarkPostAsDone(c *gin.Context) {
	uidVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var currentUserID uint
	switch v := uidVal.(type) {
	case uint:
		currentUserID = v
	case int:
		currentUserID = uint(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, currentUserID).Error; err != nil || !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Only administrators can mark posts as done"})
		return
	}

	id := c.Param("id")
	postID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Post ID"})
		return
	}

	var input MarkDoneInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := config.DB.Begin()

	var status models.Status
	if err := tx.Where("post_id = ?", postID).First(&status).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Status for this post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve status: " + err.Error()})
		return
	}

	status.Status = 0
	status.ClaimerName = input.ClaimerName
	status.ProofImage = input.ProofImage
	status.UpdatedBy = currentUserID
	status.UpdatedAt = time.Now()

	if err := tx.Save(&status).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post status: " + err.Error()})
		return
	}

	var post models.Post
	tx.First(&post, postID)

	notificationMessage := fmt.Sprintf("Laporan '%s' telah diselesaikan oleh Admin", post.Title)
	notification := models.Notification{
		PostID:  uint(postID),
		Message: notificationMessage,
		IsRead:  false,
	}

	if err := tx.Create(&notification).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification for completion: " + err.Error()})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Post marked as done successfully", "status": status})
}
