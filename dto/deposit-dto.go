package dto

type DepositDTO struct {
	ID_User     uint64 `json:"id_user" form:"id_user"`
	Amount      uint64 `json:"amount" form:"amount" binding:"required"`
	PaymentType string `json:"payment" form:"payment" binding:"required"`
	Status      uint64 `json:"status" form:"status"`
}

type DepositResponse struct {
	Id_deposit      string `json:"id_deposit" form:"id_deposit"`
	Date            string `json:"date" form:"date"`
	Id_user         uint64 `json:"id_user" form:"id_user"`
	Virtual_account string `json:"virtual_account" form:"virtual_account"`
	Url_callback    string `json:"url_callback" form:"virtual_account"`
	Amount          uint64 `json:"amount" form:"amount"`
	Status          string `json:"status" form:"status"`
}
