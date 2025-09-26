package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cyber-go/internal/handlers"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestSubmitHandler(t *testing.T) {
	// Correct way to initialize the mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Set the global DB connection in your handlers package
	handlers.SetDB(db)

	// Define the expected database operations and their results
	// The handler will likely perform a SELECT to get question data
	rows := sqlmock.NewRows([]string{"id", "weight", "selector", "options"}).
		AddRow(1, 10, "radio", "Yes,No").
		AddRow(2, 10, "checkbox", "AWS,GCP,Azure").
		AddRow(3, 15, "radio", "Yes,No")

	mock.ExpectQuery("SELECT id, weight, selector, options FROM questions").WillReturnRows(rows)

	// The handler will also perform an INSERT to save the result
	mock.ExpectExec("INSERT INTO results").WillReturnResult(sqlmock.NewResult(1, 1))

	// Create the HTTP payload
	payload := map[string]interface{}{
		"userId": "12",
		"answers": map[string]interface{}{
			"1": "Yes",
			"2": []string{"AWS", "GCP"},
			"3": "No",
		},
	}
	body, _ := json.Marshal(payload)

	// Create and execute the HTTP request
	req := httptest.NewRequest("POST", "/submit", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call the handler
	handlers.SubmitHandler(w, req)

	// Check the response
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", res.StatusCode)
	}

	// Check if all mock expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
