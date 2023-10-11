package dto

type RefundDTO struct {
	OrderID uint64  `json:"IdOrder" form:"IdOrder"`
	Amount  float64 `json:"amount" validate:"required"`
	TrxType uint64  `json:"type" form:"type" binding:"required"`
	Reason  string  `json:"reason" validate:"required"`
}
