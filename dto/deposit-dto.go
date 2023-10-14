package dto

type DepositDTO struct {
	ID_User     uint64 `json:"IdUser" form:"IdUser"`
	Ammount     uint64 `json:"ammount" form:"ammount" binding:"required"`
	PaymentType string `json:"payment" form:"payment" binding:"required"`
	Status      uint64 `json:"status" form:"status"`
}
