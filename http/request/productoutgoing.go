package request

type ProductOutgoing struct {
	Sku           string `json:"sku"`
	CountOutgoing int    `json:"count_outgoing"`
	Note          string `json:"note"`
}
