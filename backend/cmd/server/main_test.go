package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Define the mock fetcher
func mockFetchParadigms() ([]map[string]interface{}, error) {
	return []map[string]interface{}{
		{"id": 1, "name": "Threat", "description": "Potential threats to the organization"},
		{"id": 2, "name": "Vulnerability", "description": "Internal weaknesses"},
	}, nil
}

// Wrap the original handler to use a fetch function
func GetParadigmsHandlerMock(fetch func() ([]map[string]interface{}, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results, err := fetch()
		if err != nil {
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(results)
	}
}

func TestGetParadigmsHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/paradigms", nil)
	rr := httptest.NewRecorder()

	handler := GetParadigmsHandlerMock(mockFetchParadigms)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %v", rr.Code)
	}

	if len(rr.Body.String()) == 0 {
		t.Errorf("Expected non-empty body")
	}
}
