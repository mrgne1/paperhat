package main

import (
	"fmt"
	"log"
	"net/http"

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

	mux.Handle("POST /api/secrets", cfg.CreateSecretHandler())
	mux.HandleFunc("GET /api/heartbeat", handlers.Heartbeat)
	
	fmt.Printf("Serving on %v\n", port)
	log.Fatal(server.ListenAndServe())
}
