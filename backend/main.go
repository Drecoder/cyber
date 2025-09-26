package main

import (
	"log"
	"net/http"
	"time"

	"cyber-go/internal/handlers"
	"cyber-go/internal/middleware"
	"cyber-go/internal/observability" // Ensure this import path is correct
	"cyber-go/internal/util"
	"cyber-go/pkg/db"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	//1 Logger
	cleanup := util.InitLogger()
	defer cleanup()

	// 2. Tracer
	shutdown := observability.InitTracer()
	defer shutdown()

	// 3. Metrics (register + scrape)
	observability.RegisterMetrics(util.Logger)

	// 4 Inidt DB (aftrer tracer, before app start)
	myDB := db.Connect()
	defer myDB.Close()
	handlers.DB = myDB

	r := mux.NewRouter()
	r.Use(middleware.ObservabilityMiddleware(util.Logger))

	r.Handle("/metrics", promhttp.Handler())

	// REST endpoints
	r.HandleFunc("/questions", handlers.GetQuestionsHandler).Methods("GET")
	r.HandleFunc("/submit", handlers.SubmitHandler).Methods("POST")
	r.HandleFunc("/result/{userID}", handlers.ResultHandler).Methods("GET")

	// Basic HTTP server (placeholder for GraphQL)
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Cyber Service is running"))
	})
	// Rest endpoint
	r.HandleFunc("/paradigms", handlers.GetParadigmsHandler(func() ([]map[string]interface{}, error) {
		rows, err := handlers.DB.Query("SELECT * FROM paradigms")
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		var results []map[string]interface{}
		for rows.Next() {
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}
			if err := rows.Scan(valuePtrs...); err != nil {
				return nil, err
			}
			entry := make(map[string]interface{})
			for i, col := range columns {
				entry[col] = values[i]
			}
			results = append(results, entry)
		}
		return results, nil
	}))

	// GraphQL endpoint
	r.Handle("/graphql", handlers.GraphqlHandler(handlers.Schema))

	middleware.MiddlewareScraper(30 * time.Second)

	// Start server
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
