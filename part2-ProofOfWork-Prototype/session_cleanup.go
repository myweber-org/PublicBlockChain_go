
package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/db"
	"yourproject/internal/models"
)

const cleanupInterval = 24 * time.Hour

func main() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		cleanupExpiredSessions()
	}
}

func cleanupExpiredSessions() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
	result, err := db.Conn.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Failed to clean up sessions: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("Cleaned up %d expired sessions", rowsAffected)
	}
}package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/db"
	"yourproject/internal/models"
)

func main() {
	ctx := context.Background()
	database := db.GetDB()

	// Run cleanup every 24 hours
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cleanupExpiredSessions(ctx, database)
		}
	}
}

func cleanupExpiredSessions(ctx context.Context, db *db.Database) {
	cutoff := time.Now().Add(-24 * time.Hour)
	result := db.WithContext(ctx).
		Where("last_activity < ?", cutoff).
		Delete(&models.Session{})

	if result.Error != nil {
		log.Printf("Error cleaning sessions: %v", result.Error)
		return
	}

	if result.RowsAffected > 0 {
		log.Printf("Cleaned up %d expired sessions", result.RowsAffected)
	}
}