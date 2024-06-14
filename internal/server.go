package internal

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ServerInit() {
	db := dbInit(sqlite.Open("currency.db"))

	mux := http.NewServeMux()

	mux.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		getQuotationHandler(w, db)
	})

	http.ListenAndServe(":8080", mux)
}

func getQuotationHandler(w http.ResponseWriter, db *gorm.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)

	if err != nil {
		internalServerError(w, err)
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			timeoutError(w, err)
			return
		}
		internalServerError(w, err)
		return
	}
	if res.StatusCode != http.StatusOK {
		internalServerError(w, err)
		return
	}
	defer res.Body.Close()

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			timeoutError(w, err)
		} else {
			internalServerError(w, err)
		}
		return
	}

	quotationResponse := QuotationResponse{}
	err = json.Unmarshal(resBytes, &quotationResponse)
	if err != nil {
		internalServerError(w, err)
		return
	}
	db.Create(&quotationResponse.Quotation)

	quotationBytes, err := json.Marshal(&quotationResponse.Quotation)
	if err != nil {
		internalServerError(w, err)
		return
	}

	w.Write(quotationBytes)
}

func timeoutError(w http.ResponseWriter, err error) {
	log.Printf("Timeout ao buscar a cotação: %v", err)
	http.Error(w, "A busca da cotação demorou demais. Tente novamente mais tarde.", http.StatusGatewayTimeout)
}

func internalServerError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func dbInit(dialector gorm.Dialector) *gorm.DB {
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Quotation{})
	return db
}
