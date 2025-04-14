package database

import (
	"testing"
)

func TestDatabase(t *testing.T) {
	err := Connect()
	if err != nil {
		t.Fatalf("Connect: Expected no error, got %v", err)
	}

	db, err := GetDBHandle()
	if err != nil {
		t.Fatalf("GetDBHandle: Expected no error, got %v", err)
	}

	if db == nil {
		t.Fatalf("db: Expected a valid database connection, got nil")
	}

	if err := db.Ping(); err != nil {
		t.Errorf("Ping: Expected successful ping, got error: %v", err)
	}

	if err = Disconnect(); err != nil {
		t.Fatalf("Disconnect: Expected no error on disconnect, got %v", err)
	}
	db = nil
}
