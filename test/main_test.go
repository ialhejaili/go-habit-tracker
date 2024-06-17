package test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/ialhejaili/go-habit-tracker/repository"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	TestDB     *sql.DB
	TestUserID int
)

func TestMain(m *testing.M) {

	setup()

	code := m.Run()

	teardown()

	os.Exit(code)
}

func setup() {
	setupDB()
	RegisterTestUser()
}
func teardown() {
	DeleteTestUser()
	closeDB()
	log.Println("Exiting test suite")
}

func setupDB() {
	log.Println("Loading .env file...")
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dsn := os.Getenv("TEST_PG_DSN")
	if dsn == "" {
		log.Fatalf("TEST_PG_DSN is not set in the environment variables")
	} else {
		log.Println("TEST_PG_DSN is set")
	}

	log.Println("Connecting to test database...")
	TestDB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error connecting to test database: %v\n", err)
	}

	log.Println("Pinging test database...")
	err = TestDB.Ping()
	if err != nil {
		log.Fatalf("Error pinging test database: %v\n", err)
	}

	log.Println("Successfully connected to the test database")
}
func closeDB() {
	log.Println("Closing test database connection...")
	err := TestDB.Close()
	if err != nil {
		log.Fatalf("Error closing the test database: %v\n", err)
	}
}

func RegisterTestUser() {
	err := repository.RegisterUser(TestDB, "testuser", "testpassword")
	if err != nil {
		log.Fatalf("Expected no error, got %v", err)
	}

	err = TestDB.QueryRow("SELECT id FROM users WHERE username = $1", "testuser").Scan(&TestUserID)
	if err != nil {
		log.Fatalf("Expected no error, got %v", err)
	}
}

func DeleteTestUser() {
	err := repository.DeleteUser(TestDB, TestUserID)
	if err != nil {
		log.Fatalf("Expected no error, got %v", err)
	}
}
