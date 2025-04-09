package routes

import (
	"HMCTS-Developer-Challenge/database"
	"crypto/sha256"
	"database/sql"
	"fmt"
)

func LoginUser(username, password string) error {
	if username == "" || password == "" {
		return fmt.Errorf("username and password cannot be empty")
	}

	query := fmt.Sprintf("SELECT password_sha256 FROM users WHERE user = '%s'", username)
	var storedPasswordSha256 string
	err := db.GetDBHandle().QueryRow(query).Scan(&storedPasswordSha256)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user does not exist")
		} else {
			return err
		}
	}

	hashedPassword := sha256.Sum256([]byte(password))
	if storedPasswordSha256 != string(hashedPassword[:]) {
		return fmt.Errorf("invalid password")
	}

	return nil
}
