package model

import (
	"github.com/google/uuid"
)

type Bank struct {
	BankID       uuid.UUID `gorm:"type:char(36);primaryKey" json:"bank_id"`
	Name         string    `gorm:"not null" json:"name"`
	Abbreviation string    `json:"abbreviation"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	Accounts     []Account `gorm:"foreignKey:BankID" json:"accounts"`
}
