package entity

type Transaction struct {
	ID              uint64 `gorm:"primary_key:auto_increment" json:"id"`
	ID_User         uint64 `gorm:"type:int(100);index" json:"id_user"`
	User            User   `gorm:"foreignKey:ID_User" json:"-"`
	TransactionFrom uint64 `gorm:"type:varchar(255)" json:"acc_number_from"`
	TransactionTo   uint64 `gorm:"type:varchar(255)" json:"acc_number_to"`
	Date            int64  `gorm:"type:bigint" json:"date"`
	Amount          uint64 `gorm:"type:int(100)" json:"amount"`
	Status          uint64 `gorm:"type:int(100);default:1" json:"status"`
}
