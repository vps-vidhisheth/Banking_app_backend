package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	CustomerID uuid.UUID `gorm:"type:char(36);primaryKey" json:"customer_id"`
	FirstName  string    `gorm:"not null" json:"first_name"`
	LastName   string    `gorm:"not null" json:"last_name"`
	Email      string    `gorm:"unique;not null" json:"email"`
	Password   string    `gorm:"not null" json:"-"`
	Role       string    `gorm:"not null" json:"role"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	Accounts   []Account `gorm:"foreignKey:CustomerID" json:"accounts"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
