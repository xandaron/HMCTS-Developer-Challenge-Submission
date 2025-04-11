package errors

import (
	"net/http"
	"log"
	"fmt"
)

func Error(message string) error {
	return fmt.Errorf("%s", message)
}

func Errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func AddContext(err error, context string) error {
	return fmt.Errorf("%s: %s", context, err)
}

func HandleServerError(w http.ResponseWriter, err error, context string) {
	log.Printf("Error: %s\n", AddContext(err, context))
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}