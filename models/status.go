package models

import (
	"time"
)

type Status struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PostID      uint      `gorm:"uniqueIndex;not null" json:"post_id"`
	Status      int       `gorm:"default:1" json:"status"`
	ClaimerName string    `json:"claimer_name,omitempty"`
	ProofImage  string    `json:"proof_image,omitempty"`
	UpdatedBy   uint      `json:"updated_by"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Post *Post `gorm:"foreignKey:PostID" json:"-"`
}
