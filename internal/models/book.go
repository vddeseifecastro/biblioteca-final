package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	ISBN        string `gorm:"unique;not null"`
	Title       string `gorm:"not null"`
	Author      string `gorm:"not null"`
	Description string `gorm:"type:text"` // opcional
	Category    string `gorm:"not null"`
	Publisher   string
	PublishedAt *time.Time // puntero para permitir fecha vacía
	Pages       int
	Language    string `gorm:"default:'Español'"`
	Stock       int    `gorm:"default:1"`
	Available   int    `gorm:"default:1"`
	CoverImage  string
	IsActive    bool `gorm:"default:true"`
}
