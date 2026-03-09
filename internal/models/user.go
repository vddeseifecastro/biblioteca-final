package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex:idx_username_deleted,where:deleted_at IS NULL;not null"`
	Email        string `gorm:"uniqueIndex:idx_email_deleted,where:deleted_at IS NULL;not null"`
	Password     string `gorm:"not null"`
	Name         string `gorm:"not null"`
	Role         string `gorm:"default:'user'"` // Valores: "user", "bibliotecario", "admin"
	IsActive     bool   `gorm:"default:true"`
	BlockedUntil *time.Time
	Preferences  string         `gorm:"type:json"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// HashPassword genera un hash seguro para la contraseña
func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword compara la contraseña proporcionada con el hash almacenado
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
