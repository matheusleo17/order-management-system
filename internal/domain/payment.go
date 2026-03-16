package domain

type Payment struct {
	Method string `json:"method"`
	Status string `json:"status"`
}
