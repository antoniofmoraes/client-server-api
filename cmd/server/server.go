package main

import (
	"net/http"

	"github.com/antoniofmoraes/client-server-api/internal"
	"gorm.io/driver/sqlite"
)

func main() {
	db := internal.DbInit(sqlite.Open("currency.db"))

	mux := http.NewServeMux()

	mux.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		internal.GetQuotationHandler(w, db)
	})

	http.ListenAndServe(":8080", mux)
}
