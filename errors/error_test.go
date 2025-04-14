package errors

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestError(t *testing.T) {
	err := Error("test error")
	if err.Error() != "test error" {
		t.Errorf("Expected 'test error', got %v", err.Error())
	}
}

func TestErrorf(t *testing.T) {
	err := Errorf("test %s", "error")
	if err.Error() != "test error" {
		t.Errorf("Expected 'test error', got %v", err.Error())
	}
}

func TestAddContext(t *testing.T) {
	baseErr := errors.New("base error")
	contextErr := AddContext(baseErr, "context")
	if !strings.Contains(contextErr.Error(), "context: base error") {
		t.Errorf("Expected 'context: base error', got %v", contextErr.Error())
	}
}

func TestHandleServerError(t *testing.T) {
	w := httptest.NewRecorder()
	err := errors.New("test error")
	HandleServerError(w, err, "context")

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %v", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Internal Server Error") {
		t.Errorf("Expected 'Internal Server Error' in response, got %v", w.Body.String())
	}
}
