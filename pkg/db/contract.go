package db

import (
  "fmt"
)

// Contract represents a contract entity in the database
type Contract struct {
  ID          int
  ClientID    int
  Description string
  Status      string
  Code        string // Added to store contract code
}

// CreateContract adds a new contract to the database
func CreateContract(contract *Contract) (int, error) {
  var id int
  err := DB.QueryRow("INSERT INTO contracts (client_id, description, status) VALUES (?, ?, ?) RETURNING id",
    contract.ClientID, contract.Description, contract.Status).Scan(&id)
  if err != nil {
    return 0, fmt.Errorf("failed to insert contract: %v", err)
  }
  return id, nil
}

// InsertContractCode inserts contract code into the database
func InsertContractCode(id int, code string) error {
  _, err := DB.Exec("UPDATE contracts SET code = ? WHERE id = ?", code, id)
  if err != nil {
    return fmt.Errorf("failed to insert contract code: %v", err)
  }
  return nil
}

// UpdateContract updates a contract's information in the database
func UpdateContract(contract *Contract) error {
  _, err := DB.Exec("UPDATE contracts SET client_id = ?, description = ?, status = ? WHERE id = ?",
    contract.ClientID, contract.Description, contract.Status, contract.ID)
  if err != nil {
    return fmt.Errorf("failed to update contract: %v", err)
  }
  return nil
}

// UpdateContractStatus updates a contract's status in the database
func UpdateContractStatus(id int, status string) error {
  _, err := DB.Exec("UPDATE contracts SET status = ? WHERE id = ?", status, id)
  if err != nil {
    return fmt.Errorf("failed to update contract status: %v", err)
  }
  return nil
}

// DeleteContract removes a contract from the database
func DeleteContract(id int) error {
  _, err := DB.Exec("DELETE FROM contracts WHERE id = ?", id)
  if err != nil {
    return fmt.Errorf("failed to delete contract: %v", err)
  }
  return nil
}

// GetContractByID retrieves a contract from the database by ID
func GetContractByID(id int) (*Contract, error) {
  contract := &Contract{}
  err := DB.QueryRow("SELECT id, client_id, description, status, code FROM contracts WHERE id = ?", id).
    Scan(&contract.ID, &contract.ClientID, &contract.Description, &contract.Status, &contract.Code)
  if err != nil {
    return nil, fmt.Errorf("failed to get contract: %v", err)
  }
  return contract, nil
}

