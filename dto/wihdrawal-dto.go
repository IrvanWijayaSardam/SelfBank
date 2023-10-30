package dto

type WithdrawalDTO struct {
	ID_User uint64 `json:"iduser" form:"iduser"`
	Amount  uint64 `json:"amount" validate:"required"`
	To      string `json:"to" form:"to" binding:"required"`
}

type WithdrawalResponseDTO struct {
	ID     uint64 `json:"id"`
	IDUser uint64 `json:"id_user"`
	Date   string `json:"date"`
	Amount uint64 `json:"amount"`
	To     string `json:"to"`
	Status uint64 `json:"status"`
}
