package controllers

import (
	"cyber-go/internal/models"
)

// --- Scoring Logic ---
func indexOf(val string, arr []string) int {
	for i, v := range arr {
		if v == val {
			return i
		}
	}
	return 0
}

func determinePolicy(totalScore int) string {
	switch {
	case totalScore < 20:
		return "Basic Cyber Insurance"
	case totalScore < 50:
		return "Standard Cyber Insurance"
	default:
		return "Premium Cyber Insurance"
	}
}

func EvaluateAnswers(answers map[int]interface{}, questions []models.Question) (int, string) {
	totalScore := 0

	for _, q := range questions {
		ans, ok := answers[q.ID]
		if !ok {
			continue
		}
		score := 0
		switch q.Selector {
		case "radio":
			if ans.(string) == "Yes" {
				score = q.Weight
			}
		case "checkbox":
			selected := ans.([]string)
			score = q.Weight * len(selected) / len(q.Options)
		case "dropdown":
			optionIndex := indexOf(ans.(string), q.Options)
			score = q.Weight * (optionIndex + 1) / len(q.Options)
		}
		totalScore += score
	}

	policy := determinePolicy(totalScore)
	return totalScore, policy
}
