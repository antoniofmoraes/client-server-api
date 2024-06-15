package internal

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

func GetQuotationHandler(w http.ResponseWriter, db *gorm.DB) {
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

	err = InsertQuotation(db, quotationResponse)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}
		internalServerError(w, err)
		return
	}

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
