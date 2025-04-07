package repository

import (
	"database/sql"
	"film-library/internal/config"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

type Storage struct {
	db *sql.DB
}

func Connect(c config.Database) (*Storage, error) {
	const op = "storage.postgre.New"
	sqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Dbname)

	db, err := sql.Open("postgres", sqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	log.Printf("Database connected was created: %s", sqlInfo)

	dir, _ := os.Getwd()
	log.Println("Current working directory:", dir)

	if err = goose.Up(db, "./migrations"); err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return &Storage{db: db}, nil
}
