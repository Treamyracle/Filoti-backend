package models

import (
	"time"
)

type Notification struct {
	ID uint `gorm:"primaryKey" json:"id"`

	PostID    uint      `gorm:"not null;index" json:"post_id"`
	Message   string    `gorm:"not null" json:"message"`
	IsRead    bool      `gorm:"not null;default:false" json:"is_read"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	Post *Post `gorm:"foreignKey:PostID" json:"-"`
}
