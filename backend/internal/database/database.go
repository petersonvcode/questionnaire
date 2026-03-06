package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type DatabaseService struct {
	db          *sql.DB
	isConnected bool
}

var dbInstance *DatabaseService

func New() (*DatabaseService, error) {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance, nil
	}

	return &DatabaseService{isConnected: false}, nil
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *DatabaseService) Health() map[string]string {
	stats := make(map[string]string)

	if !s.isConnected {
		err := s.Connect()
		if err != nil {
			stats["status"] = "down"
			stats["error"] = fmt.Sprintf("db down: %v", err)
			slog.Error("db down", "err", err.Error())
			return stats
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		slog.Error("db down", "err", err.Error()) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *DatabaseService) Close() error {
	slog.Info("Disconnected from database")
	return s.db.Close()
}

func (s *DatabaseService) Exec(query string, args ...any) (sql.Result, error) {
	if !s.isConnected {
		err := s.Connect()
		if err != nil {
			slog.Error("failed to connect to database", "err", err.Error())
			return nil, err
		}
	}

	return s.db.Exec(query, args...)
}

func (s *DatabaseService) Query(query string, args ...any) (*sql.Rows, error) {
	if !s.isConnected {
		err := s.Connect()
		if err != nil {
			slog.Error("failed to connect to database", "err", err.Error())
			return nil, err
		}
	}

	return s.db.Query(query, args...)
}

func (s *DatabaseService) QueryRow(query string, args ...any) *sql.Row {
	if !s.isConnected {
		err := s.Connect()
		if err != nil {
			slog.Error("failed to connect to database", "err", err.Error())
			return nil
		}
	}

	return s.db.QueryRow(query, args...)
}

func getDbUrl() string {
	return os.Getenv("DATABASE_URL")
}

func (s *DatabaseService) Connect() error {
	var err error
	dbUrl := getDbUrl()
	slog.Info("Connecting to database", "url", dbUrl)
	s.db, err = sql.Open("sqlite3", dbUrl)
	if err != nil {
		slog.Error("failed to open database", "err", err.Error(), "dbUrl", dbUrl)
		return err
	}
	s.isConnected = true

	// Initializing database and migrating tables
	err = s.Migrate()
	if err != nil {
		slog.Error("failed to migrate database", "err", err.Error(), "dbUrl", dbUrl)
		return err
	}

	return nil
}

// Begins a new transaction.
// Returns the transaction, use it to execute the queries and commit or rollback the transaction.
func (s *DatabaseService) BeginTransaction() (*sql.Tx, error) {
	if !s.isConnected {
		err := s.Connect()
		if err != nil {
			slog.Error("failed to connect to database", "err", err.Error())
			return nil, err
		}
	}

	return s.db.Begin()
}
