package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/antoniofmoraes/client-server-api/internal"
)

const defaultSleepTime = 60

func main() {
	sleepTime := getSleepTime()

	if sleepTime > 0 {
		for {
			internal.GetAndSaveQuotation()
			time.Sleep(time.Duration(sleepTime) * time.Second)
		}
	} else {
		internal.GetAndSaveQuotation()
	}
}

func getSleepTime() int {
	if len(os.Args) < 2 {
		log.Printf("Tempo de intervalo automaticamente definido para %v segundos", defaultSleepTime)
		return int(defaultSleepTime)
	}

	sleepTime, err := strconv.ParseInt(os.Args[1], 0, 64)
	if err != nil {
		log.Printf("Tempo de intervalo automaticamente definido para %v segundos", defaultSleepTime)
		return int(defaultSleepTime)
	}

	return int(sleepTime)
}
