package entity

type Deposit struct {
	ID      string `gorm:"primary_key" json:"id"`
	ID_User uint64 `gorm:"type:int(100);index" json:"id_user"`
	User    User   `gorm:"foreignKey:ID_User" json:"-"`
	Date    int64  `gorm:"type:bigint" json:"date"`
	Amount  uint64 `gorm:"type:int(100)" json:"amount"`
	Status  uint64 `gorm:"type:int(100);default:1" json:"status"`
}
