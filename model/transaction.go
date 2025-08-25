package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	Deposit  = "deposit"
	Withdraw = "withdraw"
	Transfer = "transfer"
)

type Transaction struct {
	TransactionID    uuid.UUID      `gorm:"type:char(36);primaryKey" json:"transaction_id"`
	AccountID        uuid.UUID      `gorm:"type:char(36);not null" json:"account_id"`
	RelatedAccountID *uuid.UUID     `gorm:"type:char(36)" json:"related_account_id,omitempty"`
	Amount           float64        `gorm:"not null" json:"amount"`
	Type             string         `gorm:"type:enum('deposit','withdraw','transfer');not null" json:"type"`
	Note             string         `gorm:"type:varchar(255)" json:"note"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}
