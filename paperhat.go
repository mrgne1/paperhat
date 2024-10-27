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
	host := "localhost"
	hostUrl := fmt.Sprintf("%v:%v", host, port)

	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr: hostUrl,
	}

	cfg, err := handlers.NewApiConfig(hostUrl)
	if err != nil {
		log.Println(err)
		return
	}

	go cronJob(&cfg, time.Minute)

	// Website Handlers
	siteHandler := http.StripPrefix("/v1", http.FileServer(http.Dir("./site/v1")))
	mux.Handle("/v1/", siteHandler)

	// Backend Handlers
	mux.Handle("POST /api/secrets", cfg.CreateSecretHandler())
	mux.Handle("GET /api/secrets/{id}/{keyText}", cfg.ReadSecretHandler())
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
