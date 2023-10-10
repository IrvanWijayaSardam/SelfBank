package dto

type TransactionDTO struct {
	ID_User uint64 `json:"IdUser" form:"IdUser"`
	Ammount uint64 `json:"ammount" form:"ammount" binding:"required"`
	TrxType uint64 `json:"type" form:"type" binding:"required"`
	Status  uint64 `json:"status" form:"status"`
}
