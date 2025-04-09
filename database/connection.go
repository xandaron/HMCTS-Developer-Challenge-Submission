package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var db *sql.DB
var host = os.Getenv("DB_HOST")
var port = os.Getenv("DB_PORT")
var dbname = os.Getenv("DB_NAME")
var user = os.Getenv("DB_USER")
var pwd = os.Getenv("DB_PASSWORD")

func Connect() error {
	if host == "" {
		return fmt.Errorf("DB_HOST environment variable is not set")
	}
	if port == "" {
		return fmt.Errorf("DB_PORT environment variable is not set")
	}
	if dbname == "" {
		return fmt.Errorf("DB_NAME environment variable is not set")
	}
	if user == "" {
		return fmt.Errorf("DB_USER environment variable is not set")
	}
	if pwd == "" {
		return fmt.Errorf("DB_PASSWORD environment variable is not set")
	}

	var err error = nil
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pwd, host, port, dbname))
	return err
}

func Disconnect() error {
	var err error = nil
	if db.Ping() == nil {
		err = db.Close()
	}
	return err
}

func GetDBHandle() *sql.DB {
	if db.Ping() != nil {
		Connect()
	}
	return db
}
