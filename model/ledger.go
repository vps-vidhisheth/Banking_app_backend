package model

import (
	"time"

	"github.com/google/uuid"
)

type Ledger struct {
	LedgerID        uuid.UUID  `gorm:"type:char(36);primaryKey" json:"ledger_id"`
	AccountID       *uuid.UUID `gorm:"type:char(36)" json:"account_id,omitempty"`
	BankFromID      *uuid.UUID `gorm:"type:char(36)" json:"bank_from_id,omitempty"`
	BankToID        *uuid.UUID `gorm:"type:char(36)" json:"bank_to_id,omitempty"`
	Amount          float64    `gorm:"not null" json:"amount"`
	TransactionType string     `gorm:"not null" json:"transaction_type"` // deposit | withdraw | transfer
	EntryType       string     `gorm:"not null" json:"entry_type"`       // debit | credit
	Description     string     `json:"description"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
}
