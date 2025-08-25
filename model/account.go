package model

import (
	"github.com/google/uuid"
)

type Account struct {
	AccountID  uuid.UUID `gorm:"type:char(36);primaryKey" json:"account_id"`
	BankID     uuid.UUID `gorm:"type:char(36);not null" json:"bank_id"`
	CustomerID uuid.UUID `gorm:"type:char(36);not null" json:"customer_id"`
	Balance    float64   `gorm:"not null;default:1000" json:"balance"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	Ledgers    []Ledger  `gorm:"foreignKey:AccountID" json:"ledgers"`
}
