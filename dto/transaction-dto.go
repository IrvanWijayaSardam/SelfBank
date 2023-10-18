package dto

type TransactionDTO struct {
	ID_User         uint64 `json:"id_user" form:"id_user"`
	TransactionFrom uint64 `json:"acc_number_from" form:"acc_number_from"`
	TransactionTo   uint64 `json:"acc_number_to" form:"acc_number_to"`
	Amount          uint64 `json:"amount" validate:"required"`
}
