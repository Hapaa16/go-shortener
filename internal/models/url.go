package models

import "time"

type Url struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	OriginalUrl  string    `json:"original_url" gorm:"column:original_url"`
	ShortenedUrl string    `json:"shortened_url" gorm:"column:shortened_url`
	AccessCount  uint      `json:"access_count" gorm:"column:access_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
