package database

import (
	"biblioteca-final/internal/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect abre la conexión con SQLite y almacena la DB global
func Connect() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("library.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	log.Println("✅ Conectado a la base de datos SQLite")
	return nil
}

// GetDB devuelve la instancia global de la DB
func GetDB() *gorm.DB {
	return DB
}

// Migrate crea las tablas si no existen
func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Book{},
		&models.Loan{},
		&models.LoanReport{}, // NUEVO: tabla de reportes
	)
	if err != nil {
		log.Fatal("❌ Error migrando tablas:", err)
	}
	log.Println("✅ Migraciones completadas")
}
