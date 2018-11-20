package rasa

// Response is the type of response expected from the rasa server when querying
type Response struct {
	Intent        Intent        `json:"intent"`
	IntentRanking []Intent      `json:"intent_ranking"`
	Entities      []interface{} `json:"entities"`

	Text    string `json:"text"`
	Project string `json:"project"`
	Model   string `json:"model"`
}

// Intent represents a single intent object in a rasa response
type Intent struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
}
