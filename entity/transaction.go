package entity

import "time"

type Transaction struct {
	ID      uint64    `gorm:"primary_key:auto_increment" json:"id"`
	ID_User uint64    `gorm:"type:int(100);index" json:"id_user"`
	User    User      `gorm:"foreignKey:ID_User" json:"-"`
	Date    time.Time `gorm:"type:datetime;default:current_timestamp" json:"date"`
	Ammount uint64    `gorm:"type:int(100)" json:"ammount"`
	Status  uint64    `gorm:"type:int(100);default:1" json:"status"`
	TrxType uint64    `gorm:"type:int(100);check:trx_type IN (1, 2)" json:"trx_type"`
}