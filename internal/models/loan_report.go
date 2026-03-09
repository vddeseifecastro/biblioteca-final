package models

import "time"

type LoanReport struct {
	ID         uint `gorm:"primaryKey"`
	UserName   string
	BookTitle  string
	LoanDate   time.Time
	ReturnDate *time.Time
}
