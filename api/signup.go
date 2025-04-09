package api

import (
	"HMCTS-Developer-Challenge/database"
	"crypto/sha256"
	"fmt"
)

func CreateUser(username, password string) error {
	if username == "" || password == "" {
		return fmt.Errorf("username and password cannot be empty")
	}

	exists, err := CheckUserExists(username)
	if err != nil {
		return err
	} else if exists {
		return fmt.Errorf("user already exists")
	}

	query := fmt.Sprintf("INSERT INTO users (name, password_sha256) VALUES ('%s', '%s')", username, sha256.Sum256([]byte(password)))
	if _, err := db.GetDBHandle().Exec(query); err != nil {
		return err
	}
	return nil
}

func CheckUserExists(username string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(1) FROM users WHERE name = '%s'", username)
	var count int
	if err := db.GetDBHandle().QueryRow(query, username).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}
