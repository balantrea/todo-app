package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

const (
	userTable      = "users"
	todoListTable  = "todo_lists"
	listItemTable  = "list_item"
	userListsTable = "users_lists"
	todoItemsTable = "todo_items"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config, logger zerolog.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.DBName,
		cfg.Password,
		cfg.SSLMode,
	))
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		if closeErr := db.Close(); closeErr != nil {
			logger.Error().
				Err(closeErr).
				Msg("failed to close database connection")
		}

		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}
