package rasa

// Response is the type of response expected from the rasa server when querying
type Response struct {
	Intent        Intent   `json:"intent"`
	IntentRanking []Intent `json:"intent_ranking"`
	Entities      []Entity `json:"entities"`

	Text    string `json:"text"`
	Project string `json:"project"`
	Model   string `json:"model"`
}

// Intent represents a single intent object in a rasa response
type Intent struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
}

// Entity represents a RASA entity result
type Entity struct {
	Start      int     `json:"start"`
	End        int     `json:"end"`
	Value      string  `json:"value"`
	Entity     string  `json:"entity"`
	Confidence float64 `json:"confidence"`
	Extractor  string  `json:"extractor"`
}
