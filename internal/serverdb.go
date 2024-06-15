package internal

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"
)

func DbInit(dialector gorm.Dialector) *gorm.DB {
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Quotation{})
	return db
}

func InsertQuotation(db *gorm.DB, quotation QuotationResponse) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	err := db.WithContext(ctx).Create(&quotation.Quotation).Error
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Printf("Timeout ao inserir a cotação: %v", err)
			return err
		}
		return err
	}
	return nil
}
