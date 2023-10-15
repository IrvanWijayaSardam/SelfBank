package dto

type DepositDTO struct {
	ID_User     uint64 `json:"IdUser" form:"IdUser"`
	Amount      uint64 `json:"amount" form:"amount" binding:"required"`
	PaymentType string `json:"payment" form:"payment" binding:"required"`
	Status      uint64 `json:"status" form:"status"`
}
