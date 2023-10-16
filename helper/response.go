package helper

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseError struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type TransactionGroupSum struct {
	TransactionGroup string `json:"transaction_group"`
	TotalTransaction int    `json:"total_transaction"`
}

type TransactionReport struct {
	TransactionOut int `json:"transaction_out"`
	TransactionIn  int `json:"total_in"`
}

type EmptyObj struct{}

func BuildResponse(status bool, message string, data interface{}) Response {
	res := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	return res
}

func BuildErrorResponse(message string, data interface{}) ResponseError {
	res := ResponseError{
		Status:  false,
		Message: message,
	}
	return res
}
