package main

import (
	"blog/internal/db"
	"blog/internal/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Load .env if present (simple implementation without external package)
	loadEnv(".env")

	// Connect to Neon (PostgreSQL)
	db.Connect()

	// Router using stdlib mux
	mux := http.NewServeMux()

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/search", handlers.SearchPosts)
	mux.HandleFunc("/posts/new", handlers.NewPostForm)
	mux.HandleFunc("/posts", handlers.CreatePost)

	// Post routes by path pattern
	mux.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		switch {
		case strings.HasSuffix(path, "/edit") && r.Method == http.MethodGet:
			handlers.EditPostForm(w, r)
		case strings.HasSuffix(path, "/update") && r.Method == http.MethodPost:
			handlers.UpdatePost(w, r)
		case strings.HasSuffix(path, "/delete") && r.Method == http.MethodPost:
			handlers.DeletePost(w, r)
		case strings.HasSuffix(path, "/delete-row") && r.Method == http.MethodPost:
			handlers.DeletePostRow(w, r)
		case r.Method == http.MethodGet:
			handlers.ShowPost(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Inkwell running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

// loadEnv reads a simple KEY=VALUE .env file (no external dependency)
func loadEnv(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return // .env is optional
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		// Remove surrounding quotes if present
		if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
			val = val[1 : len(val)-1]
		}
		if os.Getenv(key) == "" { // don't override real env vars
			os.Setenv(key, val)
		}
	}
}
