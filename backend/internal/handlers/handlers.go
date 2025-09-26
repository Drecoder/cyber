package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt" // You need to import fmt for Sprintf
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"cyber-go/internal/controllers"
	"cyber-go/internal/models"
	"cyber-go/internal/util"

	"go.uber.org/zap"
)

var DB *sql.DB

// SetDB is a public function to set the database connection for the handlers.
func SetDB(conn *sql.DB) {
	DB = conn
}

var results = struct {
	sync.RWMutex
	data map[string]models.Result
}{data: make(map[string]models.Result)}

// ParadigmFetcher defines a fetcher function that returns paradigms from DB or mock
type ParadigmFetcher func() ([]map[string]interface{}, error)

func GetQuestionsFromDB() ([]models.Question, error) {
	rows, err := DB.Query("SELECT id, paradigm_id, text, selector, options, weight FROM questions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var q models.Question
		var paradigmID int
		var opts string
		if err := rows.Scan(&q.ID, &paradigmID, &q.Text, &q.Selector, &opts, &q.Weight); err != nil {
			return nil, err
		}
		q.Options = strings.Split(opts, ",")
		q.Paradigm = fmt.Sprintf("%d", paradigmID)
		questions = append(questions, q)
	}
	return questions, nil
}

// GetParadigmsHandler returns an HTTP handler that responds with JSON from the fetcher
func GetParadigmsHandler(fetch ParadigmFetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := fetch()
		if err != nil {
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

// GraphqlHandler returns a GraphQL HTTP handler for the provided schema
func GraphqlHandler(schema graphql.Schema) http.Handler {
	return handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})
}

func GetQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	qs, err := GetQuestionsFromDB()
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(qs)
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Answers map[int]interface{} `json:"answers"`
		UserID  string              `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Fetch all questions from DB
	rows, err := DB.Query("SELECT id, paradigm_id, text, selector, options, weight FROM questions")
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var qs []models.Question
	for rows.Next() {
		var q models.Question
		var paradigmID int
		var opts string
		if err := rows.Scan(&q.ID, &paradigmID, &q.Text, &q.Selector, &opts, &q.Weight); err != nil {
			http.Error(w, "Error scanning question", http.StatusInternalServerError)
			return
		}
		q.Options = strings.Split(opts, ",")
		q.Paradigm = fmt.Sprintf("%d", paradigmID)
		qs = append(qs, q)
	}

	// Correctly process answers from JSON decoder
	processedAnswers := make(map[int]interface{})
	for id, ans := range payload.Answers {
		if slice, ok := ans.([]interface{}); ok {
			stringSlice := make([]string, len(slice))
			for i, v := range slice {
				stringSlice[i] = v.(string)
			}
			processedAnswers[id] = stringSlice
		} else {
			processedAnswers[id] = ans
		}
	}

	totalScore, policy := controllers.EvaluateAnswers(processedAnswers, qs)

	// Convert values for DB insertion
	transactionID := uuid.New().String()
	score := int(totalScore) // or float64 if needed
	policyStr := fmt.Sprintf("%v", policy)

	_, err = DB.Exec(
		"INSERT INTO results (user_id, score, policy) VALUES (?, ?, ?)",
		payload.UserID, score, policyStr,
	)
	if err != nil {
		util.Logger.Info("Transaction %s: failed to save result for user %s — error: %v",
			zap.String("transactionID", transactionID),
			zap.String("userID", payload.UserID),
			zap.String("error", err.Error()),
		)
		http.Error(w, "Failed to save result", http.StatusInternalServerError)
		return
	}

	util.Logger.Info("Transaction %s: saved result for user %s with score %d and policy '%s'",
		zap.String("transactionID", transactionID),
		zap.String("payload.UserID", payload.UserID),
		zap.Int("score", score),
		zap.String("policyStr", policyStr),
	)

	// Store result (simulate ETL)
	results.Lock()
	results.data[payload.UserID] = models.Result{TotalScore: totalScore, Policy: policy}
	results.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Result{TotalScore: totalScore, Policy: policy})
}

func ResultHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	results.RLock()
	res, ok := results.data[userID]
	results.RUnlock()
	if !ok {
		http.Error(w, "Result not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"questions": &graphql.Field{
				Type: graphql.NewList(models.QuestionType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return GetQuestionsFromDB()
				},
			},
		},
	}),
})
