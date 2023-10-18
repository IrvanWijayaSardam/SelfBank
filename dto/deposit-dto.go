package dto

type DepositDTO struct {
	ID_User     uint64 `json:"id_user" form:"id_user"`
	Amount      uint64 `json:"amount" form:"amount" binding:"required"`
	PaymentType string `json:"payment" form:"payment" binding:"required"`
	Status      uint64 `json:"status" form:"status"`
}
