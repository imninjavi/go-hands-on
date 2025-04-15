package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	// Ambil credentials dari environment
	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") +
		"@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" +
		os.Getenv("DB_NAME")

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("DB unreachable:", err)
	}

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		appName := os.Getenv("APP_NAME")
		message := os.Getenv("WELCOME_MESSAGE")
		secretToken := os.Getenv("SECRET_TOKEN")

		response := fmt.Sprintf("App: %s\nMessage: %s\nSecretToken: %s\n", appName, message, secretToken)
		w.Write([]byte(response))
	})

	// Endpoint health check
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Insert name
	r.Post("/names", createNameHandler)

	// Get names
	r.Get("/names", getNamesHandler)

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", r)
}

type Name struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func createNameHandler(w http.ResponseWriter, r *http.Request) {
	var input Name
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO names (name) VALUES (?)", input.Name)
	if err != nil {
		http.Error(w, "Insert failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getNamesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name FROM names")
	if err != nil {
		http.Error(w, "Query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var names []Name
	for rows.Next() {
		var n Name
		if err := rows.Scan(&n.ID, &n.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		names = append(names, n)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(names)
}
