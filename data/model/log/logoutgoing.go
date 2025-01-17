package log

import (
	"ijahinventory/data/model/product"

	"time"
)

type LogOutgoing struct {
	ID            int             `gorm:"AUTO_INCREMENT" json:"id"`
	Timestamp     time.Time       `json:"timestamp"`
	ProductId     int             `json:"-"`
	Product       product.Product `gorm:"foreignkey:ProductId" json:"-"`
	SalePrice     int             `json:"sale_price"`
	TotalPrice    int             `json:"total_price"`
	CountOutgoing int             `json:"count_outgoing"`
	Note          string          `json:"note"`
}

func (LogOutgoing) TableName() string {
	return "log_outgoing"
}
