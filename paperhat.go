package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"www.github.com/mrgne1/paperhat/handlers"
)

//go:embed site/v1/*
var site embed.FS

func main() {
	fmt.Println("Starting PaperHat Server")

	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "2060"
	}

	mode := strings.ToLower(os.Getenv("MODE"))
	if mode == "" {
		mode = "standalone"
	}

	dbPath := os.Getenv("DBPATH")
	if dbPath == "" {
		dbPath = "./secrets.db"
	}

	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr: ":" + port,
	}

	cfg, err := handlers.NewApiConfig(dbPath)
	if err != nil {
		log.Println(err)
		return
	}

	go cronJob(&cfg, time.Minute)

	// Website Handlers
	if mode == "standalone" {
		fSys, err := fs.Sub(site, "site/v1")
		if err != nil {
			log.Printf("Error hosting Website: %v\n", err)
		} else {
			log.Println("Hosting Paperhat website")
			siteHandler := http.StripPrefix("/v1", http.FileServer(http.FS(fSys)))
			mux.Handle("/v1/", siteHandler)
			mux.Handle("/", http.RedirectHandler("/v1", http.StatusMovedPermanently))
		}
	}

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
