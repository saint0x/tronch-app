package db

import (
  "database/sql"
  "fmt"
  "log"
  "os"

  _ "github.com/mattn/go-sqlite3" // Import sqlite3 driver
)

var (
  DB *sql.DB
)

// Initializes the SQLite database
func Initialize(dbPath string) error {
  var err error

  // Check if database file exists
  if _, err = os.Stat(dbPath); os.IsNotExist(err) {
    // If database file doesn't exist, create it
    file, err := os.Create(dbPath)
    if err != nil {
      return fmt.Errorf("failed to create database file: %v", err)
    }
    file.Close()
  }

  // Open the database
  DB, err = sql.Open("sqlite3", dbPath)
  if err != nil {
    return fmt.Errorf("failed to open database: %v", err)
  }

  // Create tables if they don't exist
  _, err = DB.Exec(`
    CREATE TABLE IF NOT EXISTS users (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      first_name TEXT,
      last_name TEXT,
      email TEXT UNIQUE,
      password TEXT
    );

    CREATE TABLE IF NOT EXISTS clients (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      user_id INTEGER,
      name TEXT,
      email TEXT,
      FOREIGN KEY (user_id) REFERENCES users(id)
    );

    CREATE TABLE IF NOT EXISTS contracts (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      address TEXT UNIQUE,
      code TEXT,
      status TEXT,
      client_id INTEGER,
      description TEXT,
      FOREIGN KEY (client_id) REFERENCES clients(id)
    );

    CREATE TABLE IF NOT EXISTS contract_copies (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      contract_id INTEGER,
      code TEXT,
      timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
      FOREIGN KEY (contract_id) REFERENCES contracts(id)
    );
  `)
  if err != nil {
    return fmt.Errorf("failed to create tables: %v", err)
  }

  log.Println("Database initialized")
  return nil
}

// Prints the schema of the SQLite database
func PrintSchema() {
  rows, err := DB.Query("SELECT name, sql FROM sqlite_master WHERE type='table';")
  if err != nil {
    log.Fatalf("Failed to query database schema: %v", err)
  }
  defer rows.Close()

  log.Println("Database Schema:")
  for rows.Next() {
    var tableName, tableSchema string
    if err := rows.Scan(&tableName, &tableSchema); err != nil {
      log.Fatalf("Failed to scan row: %v", err)
    }
    log.Printf("Table: %s\nSchema: %s\n", tableName, tableSchema)
  }
}

// Close closes the database connection
func Close() {
  if DB != nil {
    DB.Close()
  }
}
