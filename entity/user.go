package entity

type User struct {
	ID            uint64 `gorm:"primary_key:auto_increment" json:"id"`
	Namadepan     string `gorm:"type:varchar(255)" json:"nama_depan"`
	Namabelakang  string `gorm:"type:varchar(255)" json:"nama_belakang"`
	Email         string `gorm:"type:varchar(255)" json:"email"`
	Username      string `gorm:"type:varchar(255)" json:"username"`
	Password      string `gorm:"->;<-;not null" json:"-"`
	Telephone     string `gorm:"type:varchar(255)" json:"telp"`
	Jk            string `gorm:"type:varchar(255)" json:"jk"`
	Profile       string `gorm:"type:varchar(255)" json:"profile"`
	Token         string `gorm:"-" json:"token,omitempty"`
	Balance       string `gorm:"-" json:"balance,omitempty"`
	AccountNumber uint64 `gorm:"type:varchar(255)" json:"acc_number"`
	IdRole        uint64 `gorm:"type:bigint" json:"idrole"`
	Status        uint64 `gorm:"type:int(100);default:1" json:"status"`
	IsVerified    bool   `gorm:"type:boolean" json:"is_verified`
}
