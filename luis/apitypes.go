/*
	apitypes.go

	This file contains the types for the LUIS API
*/
package luis

type LUISResponse struct {
	Query             string            `json:"query"`
	TopScoringIntent  Intent            `json:"topScoringIntent"`
	Intents           []Intent          `json:"intents"`
	Entities          []Entity          `json:"entities"`
	SentimentAnalysis SentimentAnalysis `json:"sentimentAnalysis"`
}

type Intent struct {
	Intent string  `json:"intent"`
	Score  float64 `json:"score"`
}

type Entity struct {
	Entity     string  `json:"entity"`
	Type       string  `json:"type"`
	StartIndex int     `json:"startIndex"`
	EndIndex   int     `json:"endIndex"`
	Score      float64 `json:"score"`
	Resolution struct {
		Value  string   `json:"value"`
		Unit   string   `json:"unit"`
		Values []string `json:"values"`
	} `json:"resolution"`
}

type SentimentAnalysis struct {
	Label string  `json:"string"`
	Score float64 `json:"score"`
}
