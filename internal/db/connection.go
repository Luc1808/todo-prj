package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	name := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")

	connStr := fmt.Sprintf("host=%s user=%s port=%s dbname=%s password=%s  sslmode=disable", host, user, port, name, password)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		panic("Could not start DB connection.")
	}
	// defer DB.Close()

	createTables()
}

func createTables() {
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(45) UNIQUE NOT NULL,
		password TEXT
	);`)
	if err != nil {
		log.Fatalf("Could not create users table: %v", err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS todo (
		id SERIAL PRIMARY KEY,
		title VARCHAR(45),
		description TEXT,
		complete BOOLEAN,
		priority VARCHAR(20),
		category VARCHAR(20),
		createdAt timestamp,
		dueDate timestamp,
		userID INTEGER,
		FOREIGN KEY (userID) REFERENCES users(id)
	);`)
	if err != nil {
		panic("Could not create todo table.")
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS token (
		id SERIAL PRIMARY KEY,
		token TEXT,
		userID INTEGER,
		createdAt timestamp,
		expiresAt timestamp,
		FOREIGN KEY (userID) REFERENCES users(id)
	);`)
	if err != nil {
		panic("Could not create token table.")
	}
}
