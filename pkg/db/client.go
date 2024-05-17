package db

import (
	"database/sql"
	"fmt"
)

// Client entity in the database
type Client struct {
	ID      int
	Name    string
	Email   string
	Address string
}

// Adds a new client to the database
func CreateClient(db *sql.DB, client *Client) error {
	_, err := db.Exec("INSERT INTO clients (name, email, address) VALUES (?, ?, ?)",
		client.Name, client.Email, client.Address)
	if err != nil {
		return fmt.Errorf("failed to insert client: %v", err)
	}
	return nil
}

// Retrieves a client from the database by ID
func GetClientByID(db *sql.DB, id int) (*Client, error) {
	client := &Client{}
	err := db.QueryRow("SELECT id, name, email, address FROM clients WHERE id = ?", id).
		Scan(&client.ID, &client.Name, &client.Email, &client.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %v", err)
	}
	return client, nil
}

// Updates a client's information in the database
func UpdateClient(db *sql.DB, client *Client) error {
	_, err := db.Exec("UPDATE clients SET name = ?, email = ?, address = ? WHERE id = ?",
		client.Name, client.Email, client.Address, client.ID)
	if err != nil {
		return fmt.Errorf("failed to update client: %v", err)
	}
	return nil
}

// Removes a client from the database
func DeleteClient(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM clients WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete client: %v", err)
	}
	return nil
}
