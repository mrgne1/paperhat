package handlers

import (
	"io"
	"log"
	"net/http"
)

func (c *ApiConfig) CreateSecretHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("CreateSecretHandler Error: %v\n", err)
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(500)
			w.Write([]byte("Error reading secret from HTTP message"))
			return
		}

		log.Printf("CreateSecretHandler Body: '%v'\n", string(body))
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("Secret Created"))
	})
}
