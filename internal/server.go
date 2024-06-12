package internal

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

func Server() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /cotacao", getCurrenncyHandler)

	http.ListenAndServe(":8080", mux)
}

func getCurrenncyHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://httpbin.org/delay/5", nil)
	// req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
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

	w.Header().Set("Content-Type", "application/json")
	w.Write(resBytes)
	w.WriteHeader(http.StatusOK)
}

func timeoutError(w http.ResponseWriter, err error) {
	log.Printf("Timeout ao buscar a cotação: %v", err)
	http.Error(w, "A busca da cotação demorou demais. Tente novamente mais tarde.", http.StatusGatewayTimeout)
}

func internalServerError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
