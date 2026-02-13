package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func MustGetPostgresPath() string {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
}
