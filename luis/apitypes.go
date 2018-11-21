package luis

// Response represents the json response from a LUIS query request
type Response struct {
	Query             string            `json:"query"`
	TopScoringIntent  Intent            `json:"topScoringIntent"`
	Intents           []Intent          `json:"intents"`
	Entities          []Entity          `json:"entities"`
	SentimentAnalysis SentimentAnalysis `json:"sentimentAnalysis"`
}

// Intent represents a matched LUIS intent
type Intent struct {
	Intent string  `json:"intent"`
	Score  float64 `json:"score"`
}

// Entity represents a LUIS entity
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

// SentimentAnalysis will show up on LUIS responses
// if it has been enabled
type SentimentAnalysis struct {
	Label string  `json:"string"`
	Score float64 `json:"score"`
}
