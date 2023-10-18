package dto

type RefundDTO struct {
	OrderID uint64  `json:"idorder" form:"idorder"`
	Amount  float64 `json:"amount" validate:"required"`
	TrxType uint64  `json:"type" form:"type" binding:"required"`
	Reason  string  `json:"reason" validate:"required"`
}
