package db

import (
	"time"

	"gorm.io/gorm"
)

type Submission struct {
	ID            string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ChallengeID   string         `gorm:"not null;index"`
	UserID        string         `gorm:"not null;index"`
	Language      string         `gorm:"not null"`
	Code          string         `gorm:"type:text;not null"`
	Stdout        string         `gorm:"type:text"`
	Stderr        string         `gorm:"type:text"`
	CompileOutput string         `gorm:"type:text"`
	Message       string         `gorm:"type:text"`
	StatusCode    int            `gorm:"not null;default:0"`
	StatusDesc    string         `gorm:"not null;default:'Pending'"`
	Memory        int            `gorm:"not null;default:0"`
	Time          float64        `gorm:"not null;default:0"`
	Token         string         `gorm:"type:varchar(255)"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
