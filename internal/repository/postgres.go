package repository

import (
	"context"
	"database/sql"
	"film-library/internal/config"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) DB() *sql.DB {
	return s.db
}

func Connect(cfg config.Database) (*Storage, error) {
	const op = "storage.postgre.New"

	// Формируем строку подключения с актуальными значениями
	sqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", cfg.Host),
		getEnv("DB_PORT", cfg.Port),
		getEnv("DB_USER", cfg.User),
		getEnv("DB_PASSWORD", cfg.Password),
		getEnv("DB_NAME", cfg.Dbname),
	)

	log.Printf("Attempting to connect to database with: %s", hidePassword(sqlInfo))

	db, err := sql.Open("postgres", sqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем соединение с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to ping database: %w", op, err)
	}

	log.Println("Database connection established")

	// Получаем абсолютный путь к миграциям
	migrationsDir, err := filepath.Abs("./migrations")
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get migrations path: %w", op, err)
	}

	log.Printf("Applying migrations from: %s", migrationsDir)
	if err := goose.Up(db, migrationsDir); err != nil {
		return nil, fmt.Errorf("%s: failed to apply migrations: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// hidePassword скрывает пароль в логах
func hidePassword(connStr string) string {
	const passwordKey = "password="
	start := strings.Index(connStr, passwordKey)
	if start == -1 {
		return connStr
	}
	start += len(passwordKey)
	end := strings.Index(connStr[start:], " ")
	if end == -1 {
		return connStr[:start] + "***"
	}
	return connStr[:start] + "***" + connStr[start+end:]
}
