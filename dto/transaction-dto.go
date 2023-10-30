package dto

type TransactionDTO struct {
	ID_User         uint64 `json:"id_user" form:"id_user"`
	TransactionFrom uint64 `json:"acc_number_from" form:"acc_number_from"`
	TransactionTo   uint64 `json:"acc_number_to" form:"acc_number_to" validate:"required"`
	Amount          uint64 `json:"amount" validate:"required"`
}

type TransactionResponse struct {
	ID                uint64 `json:"id"`
	IDUser            uint64 `json:"id_user"`
	AccountNumberFrom uint64 `json:"acc_number_from"`
	AccountNumberTo   uint64 `json:"acc_number_to"`
	Date              string `json:"date"`
	Amount            uint64 `json:"amount"`
	Status            uint64 `json:"status"`
}
