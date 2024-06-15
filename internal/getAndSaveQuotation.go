package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type GetQuotationResponse struct {
	Bid float32 `json:"bid"`
}

const errorMessage = "Erro: falha ao buscar cotacao"

func GetAndSaveQuotation() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		logError(err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Erro: timeout ao buscar cotacao")
			return
		}
		logError(err)
		return
	}
	if res.StatusCode != http.StatusOK {
		logError(err)
		return
	}
	defer res.Body.Close()

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		logError(err)
		return
	}

	var quotation GetQuotationResponse
	err = json.Unmarshal(resBytes, &quotation)
	if err != nil {
		logError(err)
		return
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		logError(err)
		return
	}
	defer file.Close()

	_, err = file.Write([]byte(fmt.Sprintf("Dólar: %v", quotation.Bid)))
	if err != nil {
		logError(err)
		return
	}

	log.Printf("Cotação gravada com sucesso: %v", quotation.Bid)
}

func logError(err error) {
	log.Println(errorMessage)
	log.Println(err.Error())
}
