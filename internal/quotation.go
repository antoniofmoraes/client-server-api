package internal

import (
	"gorm.io/gorm"
)

type Quotation struct {
	gorm.Model `json:"-"`
	Code       string  `json:"code"`
	CreateDate string  `json:"create_date" gorm:"primaryKey;unique;not null"`
	Bid        Float32 `json:"bid"`
}

type QuotationResponse struct {
	Quotation Quotation `json:"USDBRL"`
}
