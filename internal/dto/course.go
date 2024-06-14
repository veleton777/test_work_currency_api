package dto

type ConvertCurrencyReq struct {
	From   string  `json:"from" validate:"required" example:"USD"`
	To     string  `json:"to" validate:"required" example:"BTC"`
	Amount float64 `json:"amount" validate:"required,gt=0" example:"1"`
}

type ConvertCurrencyResp struct {
	Course float64 `json:"course" example:"0.00001444655"`
}
