package errors

import (
	"net/http"
	"log"
	"fmt"
)

func NewBaseError(message string) error {
	return fmt.Errorf("%s", message)
}

func New(err error, context string) error {
	return fmt.Errorf("%s: %s", context, err)
}

func HandleServerError(w http.ResponseWriter, err error, context string) {
	log.Printf("Error: %s\n", New(err, context))
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}