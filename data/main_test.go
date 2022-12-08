package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var testQueries *PostgresStore

func TestMain(m *testing.M) {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	postgresqlDbInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_SSLMODE"))

	db, err := sql.Open("postgres", postgresqlDbInfo)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = &PostgresStore{db: db}

	os.Exit(m.Run())
}
