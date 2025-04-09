package db

import (
	"fmt"
)

func GetUserID(username string) (uint, error) {
	query := fmt.Sprintf("SELECT id FROM users WHERE name = '%s'", username)
	dbHandle, err := GetDBHandle()
	if err != nil {
		return 0, err
	}

	var userID uint
	if err := dbHandle.QueryRow(query).Scan(&userID); err != nil {
		return 0, err
	}
	return userID, nil
}

func GetUserName(userID uint) (string, error) {
	query := fmt.Sprintf("SELECT name FROM users WHERE id = %v", userID)
	dbHandle, err := GetDBHandle()
	if err != nil {
		return "", err
	}

	var username string
	if err := dbHandle.QueryRow(query).Scan(&username); err != nil {
		return "", err
	}
	return username, nil
}

func GetUserPassword(userID uint) (string, error) {
	query := fmt.Sprintf("SELECT password_sha256 FROM users WHERE id = %v", userID)
	dbHandle, err := GetDBHandle()
	if err != nil {
		return "", err
	}

	var password string
	if err := dbHandle.QueryRow(query).Scan(&password); err != nil {
		return "", err
	}
	return password, nil
}