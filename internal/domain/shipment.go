package domain

type Shipment struct {
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
}
