package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"www.github.com/mrgne1/paperhat/handlers"
)

func main() {
	fmt.Println("Starting PaperHat Server")

	port := "2060"

	mux := http.NewServeMux()

	server := http.Server{
		Handler: mux,
		Addr: ":" + port,
	}

	cfg, err := handlers.NewApiConfig()
	if err != nil {
		log.Println(err)
		return
	}

	go cronJob(&cfg, time.Minute)

	mux.Handle("POST /api/secrets", cfg.CreateSecretHandler())
	mux.Handle("GET /api/secrets/{id}", cfg.ReadSecretHandler())
	mux.HandleFunc("GET /api/heartbeat", handlers.Heartbeat)
	
	fmt.Printf("Serving on %v\n", port)
	log.Fatal(server.ListenAndServe())
}

func cronJob(cfg *handlers.ApiConfig, t time.Duration) {
	for true {
		cfg.DeleteExpired()
		time.Sleep(t)
	}
}
