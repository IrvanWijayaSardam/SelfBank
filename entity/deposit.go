package entity

import "time"

type Deposit struct {
	ID      uint64    `gorm:"primary_key:auto_increment" json:"id"`
	ID_User uint64    `gorm:"type:int(100);index" json:"id_user"`
	User    User      `gorm:"foreignKey:ID_User" json:"-"`
	Date    time.Time `gorm:"type:datetime;default:current_timestamp" json:"date"`
	Amount  uint64    `gorm:"type:int(100)" json:"amount"`
	Status  uint64    `gorm:"type:int(100);default:1" json:"status"`
}
