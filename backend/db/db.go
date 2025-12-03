package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func Init() {
	var err error
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database: ", err)
	}

	log.Println("Connected to database successfully")
}

