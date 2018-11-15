/*
	apitypes.go

	This file contains the types for the LUIS API
*/
package luis

type LUISResponse struct {
	Query            string           `json:"query"`
	TopScoringIntent TopScoringIntent `json:"topScoringIntent"`
	Entities         []Entity         `json:"entities"`
}

type TopScoringIntent struct {
	Intent string  `json:"intent"`
	Score  float64 `json:"score"`
}

type Entity struct {
	Entity     string `json:"entity"`
	Type       string `json:"type"`
	StartIndex int    `json:"startIndex"`
	EndIndex   int    `json:"endIndex"`
	Resolution struct {
		Value string `json:"value"`
	} `json:"resolution"`
}
