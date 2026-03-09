package models

import (
	"time"

	"gorm.io/gorm"
)

type Loan struct {
	gorm.Model
	UserID    uint `gorm:"not null"`
	BookID    uint `gorm:"not null"`
	DueDate   *time.Time
	Returned  bool `gorm:"default:false"`
	Book      Book
	User      User
	IsOverdue bool `gorm:"-"`
}
