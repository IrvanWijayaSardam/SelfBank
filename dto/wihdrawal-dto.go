package dto

type WithdrawalDTO struct {
	ID_User uint64 `json:"iduser" form:"iduser"`
	Amount  uint64 `json:"amount" validate:"required"`
	To      string `json:"to" form:"to" binding:"required"`
}
