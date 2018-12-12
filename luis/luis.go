/*
Package luis is a luis middleware for gobbl

This file contains methods for accessing the Microsoft LUIS API
This middleware will take the c.Request.Text property and use it to query a LUIS endpoint
The top scoring intent will be stored in the "intent" flag
The entire result body will be stored in the "luis" flag */
package luis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/calebhiebert/gobbl"
)

// LUIS is a LUIS api object
type LUIS struct {
	client        *http.Client
	minConfidence float64
	endpoint      string
}

// New creates a new LUIS instance, this stores the endpoint for calling
func New(endpoint string) (*LUIS, error) {
	parsedEndpoint, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	q := parsedEndpoint.Query()

	delete(q, "q")

	queryString := "?"

	for k, v := range q {
		if len(queryString) > 1 {
			queryString += "&"
		}

		queryString += url.QueryEscape(k)
		queryString += "="
		queryString += url.QueryEscape(v[0])
	}

	if _, exists := q["subscription-key"]; !exists {
		return nil, errors.New("Missing Subscription Key")
	}

	config := LUIS{
		endpoint:      fmt.Sprintf("https://%s%s%s", parsedEndpoint.Host, parsedEndpoint.Path, queryString),
		minConfidence: 0.65,
		client: &http.Client{
			Timeout: 6 * time.Second,
		},
	}

	return &config, nil
}

// Middleware returns the LUIS middleware that will query luis with the Text property from the generic request
func Middleware(luis *LUIS) gbl.MiddlewareFunction {
	return func(c *gbl.Context) {

		if c.Request.Text == "" {
			c.Next()
			return
		}

		response, err := luis.Query(c.Request.Text)
		if err != nil {
			c.Error(fmt.Sprintf("LUIS Error %v", err))
			c.Next()
			return
		}

		if response.TopScoringIntent.Intent != "" && response.TopScoringIntent.Score >= luis.minConfidence {
			c.Flag("intent", strings.TrimSpace(response.TopScoringIntent.Intent))
		}

		c.Flag("luis", response)

		entities := make(map[string][]string)

		for _, entity := range response.Entities {
			if _, ok := entities[entity.Type]; !ok {
				entities[entity.Type] = []string{}
			}

			if entity.Resolution.Values != nil {
				entities[entity.Type] = append(entities[entity.Type], entity.Resolution.Values...)
			} else if entity.Resolution.Value != "" {
				entities[entity.Type] = append(entities[entity.Type], entity.Resolution.Value)
			} else if strings.TrimSpace(entity.Entity) != "" {
				entities[entity.Type] = append(entities[entity.Type], entity.Entity)
			}

			c.Log(50, fmt.Sprintf("Processing LUIS Entity %v", entity), "LUIS")
		}

		for entityType, entityValues := range entities {
			c.Flag("luis:e:"+entityType, entityValues)
		}

		c.Next()
	}
}

// Query will make a query against the LUIS api
func (l LUIS) Query(queryString string) (*Response, error) {
	if len(queryString) > 500 {
		queryString = string([]rune(queryString)[:500])
	}

	resp, err := l.client.Get(l.endpoint + "&q=" + url.QueryEscape(queryString))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var luisError interface{}

		err = json.Unmarshal(body, &luisError)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("LUIS Error %+v", luisError)
	}

	luisResponse := &Response{}

	err = json.Unmarshal(body, luisResponse)
	if err != nil {
		return nil, err
	}

	return luisResponse, nil
}
