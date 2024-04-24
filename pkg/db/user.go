package db

import (
  "fmt"
  "log"
)

// User represents a user in the database
type User struct {
  ID        int
  FirstName string
  LastName  string
  Email     string
  Password  string
}

// CreateUser adds a new user to the database
func CreateUser(user User) error {
  _, err := DB.Exec("INSERT INTO users (first_name, last_name, email, password) VALUES (?, ?, ?, ?)",
    user.FirstName, user.LastName, user.Email, user.Password)
  if err != nil {
    return fmt.Errorf("failed to create user: %v", err)
  }
  log.Printf("User %s added", user.Email)
  return nil
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(email string) (User, error) {
  var user User
  err := DB.QueryRow("SELECT id, first_name, last_name, email, password FROM users WHERE email = ?", email).
    Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password)
  if err != nil {
    return User{}, fmt.Errorf("failed to get user: %v", err)
  }
  return user, nil
}

// UpdateUser updates a user's details
func UpdateUser(user User) error {
  _, err := DB.Exec("UPDATE users SET first_name = ?, last_name = ?, email = ?, password = ? WHERE id = ?",
    user.FirstName, user.LastName, user.Email, user.Password, user.ID)
  if err != nil {
    return fmt.Errorf("failed to update user: %v", err)
  }
  log.Printf("User %s updated", user.Email)
  return nil
}

// DeleteUser deletes a user by ID
func DeleteUser(userID int) error {
  _, err := DB.Exec("DELETE FROM users WHERE id = ?", userID)
  if err != nil {
    return fmt.Errorf("failed to delete user: %v", err)
  }
  log.Printf("User with ID %d deleted", userID)
  return nil
}
