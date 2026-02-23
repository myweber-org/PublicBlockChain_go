package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/db"
	"yourproject/internal/models"
)

func main() {
	ctx := context.Background()
	database, err := db.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	for {
		now := time.Now()
		result, err := database.ExecContext(ctx,
			"DELETE FROM user_sessions WHERE expires_at < ?",
			now,
		)
		if err != nil {
			log.Printf("Error cleaning sessions: %v", err)
		} else {
			rows, _ := result.RowsAffected()
			log.Printf("Cleaned %d expired sessions", rows)
		}

		time.Sleep(24 * time.Hour)
	}
}