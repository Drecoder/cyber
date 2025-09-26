package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"cyber-go/internal/handlers"
)

func TestGetQuestionsHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	handlers.SetDB(db)

	// Corrected: Add the missing "paradigm_id" column
	rows := sqlmock.NewRows([]string{"id", "paradigm_id", "text", "selector", "options", "weight"}).
		AddRow(1, 101, "Question 1", "radio", "Yes,No", 10)

	// Corrected: The mock query must also include "paradigm_id"
	mock.ExpectQuery("SELECT id, paradigm_id, text, selector, options, weight FROM questions").
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/questions", nil)
	w := httptest.NewRecorder()

	handlers.GetQuestionsHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", res.StatusCode)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSubmitHandler(t *testing.T) {
	// 1. Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	handlers.SetDB(db)

	// 2. Define expected database interactions
	// This query must match the one in your handler exactly
	rows := sqlmock.NewRows([]string{"id", "paradigm_id", "text", "selector", "options", "weight"}).
		AddRow(1, 101, "Question 1", "radio", "Yes,No", 10).
		AddRow(2, 102, "Question 2", "checkbox", "AWS,GCP,Azure", 10).
		AddRow(3, 103, "Question 3", "radio", "Yes,No", 15)

	// Mocks the database call that fetches all questions for evaluation.
	mock.ExpectQuery("SELECT id, paradigm_id, text, selector, options, weight FROM questions").
		WillReturnRows(rows)

	// Mocks the database call that saves the result.
	//mock.ExpectExec("INSERT INTO results").WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO results").
		WithArgs("12", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	// 3. Create and execute the HTTP request
	// Correct payload format using a map for "answers"
	payload := map[string]any{
		"userId": "12",
		"answers": map[string]any{
			"1": "Yes",
			"2": []string{"AWS"},
			"3": "No",
		},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/submit", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handlers.SubmitHandler(w, req)

	// 4. Check the response
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", res.StatusCode)
	}

	// 5. Ensure all mock expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
