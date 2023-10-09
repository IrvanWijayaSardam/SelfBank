package dto

type RegisterDTO struct {
	Namadepan    string `json:"namadepan" form:"namadepan" binding:"required"`
	Namabelakang string `json:"namabelakang" form:"namabelakang" binding:"required"`
	Email        string `json:"email" form:"email" binding:"required,email"`
	Username     string `json:"username" form:"username" binding:"required"`
	Password     string `json:"password" form:"password" binding:"required"`
	Telephone    string `json:"telp" form:"telp" binding:"required"`
	Jk           string `json:"jk" form:"jk" binding:"required"`
	Status       uint64 `json:"status" form:"status"`
	IdRole       uint64 `json:"idrole" form:"idrole" binding:"required"`
}
