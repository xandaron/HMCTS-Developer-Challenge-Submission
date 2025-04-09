package db

import (
	"fmt"
	"os"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Connect() error {
	dbname := os.Getenv("MYSQL_DATABASE")
	if dbname == "" {
		return fmt.Errorf("MYSQL_DATABASE environment variable is not set")
	}

	user := os.Getenv("MYSQL_USER")
	if user == "" {
		return fmt.Errorf("MYSQL_USER environment variable is not set")
	}

	pwd := os.Getenv("MYSQL_PASSWORD")
	if pwd == "" {
		return fmt.Errorf("MYSQL_PASSWORD environment variable is not set")
	}

	var err error = nil
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(http://mysql-db:3306)/%s", user, pwd, dbname))
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
