package models

import "github.com/graphql-go/graphql"

var QuestionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Question",
	Fields: graphql.Fields{
		"id":         &graphql.Field{Type: graphql.Int},
		"paradigmId": &graphql.Field{Type: graphql.Int},
		"text":       &graphql.Field{Type: graphql.String},
		"selector":   &graphql.Field{Type: graphql.String},
		"options":    &graphql.Field{Type: graphql.NewList(graphql.String)},
		"weight":     &graphql.Field{Type: graphql.Int},
	},
})

type Paradigm struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Question struct {
	ID       int      `json:"id"`
	Paradigm string   `json:"paradigm"`
	Text     string   `json:"text"`
	Weight   int      `json:"weight"`
	Selector string   `json:"selector"`
	Options  []string `json:"options"`
}

type Answer struct {
	QuestionID int         `json:"questionId"`
	Response   interface{} `json:"response"`
}

type Result struct {
	TotalScore int    `json:"totalScore"`
	Policy     string `json:"policy"`
}
