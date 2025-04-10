package errors

import (
	"net/http"
	"log"
)

func HandleServerError(w http.ResponseWriter, err error, context string) {
	log.Printf("Error in %s: %s\n", context, err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}