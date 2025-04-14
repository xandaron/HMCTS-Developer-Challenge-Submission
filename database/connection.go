package database

import (
	"HMCTS-Developer-Challenge/errors"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"fmt"
)

var db *sql.DB
var host = os.Getenv("DB_HOST")
var port = os.Getenv("DB_PORT")
var dbname = os.Getenv("DB_NAME")
var user = os.Getenv("DB_USER")
var pwd = os.Getenv("DB_PASSWORD")

func Connect() error {
	if host == "" {
		return errors.Error("DB_HOST environment variable is not set")
	}
	if port == "" {
		return errors.Error("DB_PORT environment variable is not set")
	}
	if dbname == "" {
		return errors.Error("DB_NAME environment variable is not set")
	}
	if user == "" {
		return errors.Error("DB_USER environment variable is not set")
	}
	if pwd == "" {
		return errors.Error("DB_PASSWORD environment variable is not set")
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
	db = nil
	return err
}

func GetDBHandle() (*sql.DB, error) {
	var err error = nil
	if db.Ping() != nil {
		err = Connect()
	}
	return db, err
}
