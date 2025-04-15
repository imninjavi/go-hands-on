package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		appName := os.Getenv("APP_NAME")
		message := os.Getenv("WELCOME_MESSAGE")
		secretToken := os.Getenv("SECRET_TOKEN")

		response := fmt.Sprintf("App: %s\nMessage: %s\nSecretToken: %s\n", appName, message, secretToken)
		w.Write([]byte(response))
	})

	// Health check endpoint
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
