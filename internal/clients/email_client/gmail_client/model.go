package gmailclient

import (
	"paytm-project/internal/models"
	"time"
)

type TokenStorage struct {
	models.BaseModel
	UserID       string    `gorm:"not null"` // Your user identifier
	AccessToken  string    `gorm:"not null"`
	RefreshToken string    `gorm:"not null"`
	TokenType    string    `gorm:"default:'Bearer'"`
	Expiry       time.Time `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
