package domain

type Discount struct {
	Code   string  `json:"code"`
	Amount float64 `json:"amount"`
}
