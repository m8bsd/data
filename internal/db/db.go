package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("✅ Connected to Neon database")
	createSchema()
}

func createSchema() {
	query := `
	CREATE TABLE IF NOT EXISTS posts (
		id         SERIAL PRIMARY KEY,
		title      VARCHAR(255) NOT NULL,
		slug       VARCHAR(255) UNIQUE NOT NULL,
		excerpt    TEXT,
		content    TEXT NOT NULL,
		author     VARCHAR(100) NOT NULL DEFAULT 'Admin',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug);
	CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);
	`

	if _, err := DB.Exec(query); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}
	fmt.Println("✅ Schema ready")
}
