package entity

import "gorm.io/gorm"

type PaymentToken struct {
	gorm.Model
	DepositID    string `gorm:"index;not null"`
	PaymentToken string `gorm:"type:varchar(255);not null"`
	VirtualAcc   string `gorm:"type:varchar(255);not null"`
	CallbackUrl  string `gorm:"type:varchar(255);not null"`
}
