package db

func GetUserID(username string) (uint, error) {
	dbHandle, err := GetDBHandle()
	if err != nil {
		return 0, err
	}

	var userID uint
	if err := dbHandle.QueryRow("SELECT id FROM users WHERE name = ?", username).Scan(&userID); err != nil {
		return 0, err
	}
	return userID, nil
}

func GetUserName(userID uint) (string, error) {
	dbHandle, err := GetDBHandle()
	if err != nil {
		return "", err
	}

	var username string
	if err := dbHandle.QueryRow("SELECT name FROM users WHERE id = ?", userID).Scan(&username); err != nil {
		return "", err
	}
	return username, nil
}

func GetUserPassword(userID uint) (string, error) {
	dbHandle, err := GetDBHandle()
	if err != nil {
		return "", err
	}

	var password string
	if err := dbHandle.QueryRow("SELECT password_sha256 FROM users WHERE id = ?", userID).Scan(&password); err != nil {
		return "", err
	}
	return password, nil
}
