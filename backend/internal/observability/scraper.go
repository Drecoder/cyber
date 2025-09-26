package observability

import (
	"net/http"
	// "strconv"
	// "strings"
	"time"

	"cyber-go/internal/util"

	"go.uber.org/zap"
)

// ScrapeMetrics fetches /metrics endpoint and logs specific metrics
func ScrapeMetrics() {
	resp, err := http.Get("http://localhost:8080/metrics")
	if err != nil {
		util.Logger.Error("Error fetching metrics", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	util.Logger.Info("Scraper tick: metrics endpoint available")
}

// parseAndLogMetric parses a metric line and logs it
// func parseAndLogMetric(text, metricName string) {
// 	lines := strings.Split(text, "\n")
// 	for _, line := range lines {
// 		if strings.HasPrefix(line, metricName) {
// 			parts := strings.Fields(line)
// 			if len(parts) > 1 {
// 				value, err := strconv.ParseFloat(parts[len(parts)-1], 64)
// 				if err == nil {
// 					util.Logger.Info("Metric value",
// 						zap.String("metric", metricName),
// 						zap.Float64("value", value),
// 					)
// 					return
// 				}
// 			}
// 			util.Logger.Warn("Metric found but could not parse value", zap.String("line", line))
// 			return
// 		}
// 	}
// 	util.Logger.Warn("Metric not found", zap.String("metric", metricName))
// }

// StartScraper runs the scraper periodically
func StartScraper(interval time.Duration) {
	go func() {
		time.Sleep(3 * time.Second) // wait for server to start
		for {
			ScrapeMetrics()
			time.Sleep(interval)
		}
	}()
}
