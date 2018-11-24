// Package rasa contains a rasa middleware to extract nlp data
package rasa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/calebhiebert/gobbl"
)

// API impiments the RASA api
type API struct {
	client   *http.Client
	endpoint string
}

// New creates a new RASA instance, this stores the endpoint for calling
func New(endpoint string) (*API, error) {
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

	protocol := "http:"

	if strings.HasPrefix(endpoint, "https:") {
		protocol = "https:"
	}

	config := API{
		endpoint: fmt.Sprintf("%s//%s/parse%s", protocol, parsedEndpoint.Host, queryString),
		client: &http.Client{
			Timeout: 6 * time.Second,
		},
	}

	return &config, nil
}

// Middleware returns the RASA middleware that will query RASA with the Text property from the generic request
func Middleware(rasa *API) gbl.MiddlewareFunction {
	return func(c *gbl.Context) {

		if c.Request.Text == "" {
			c.Next()
			return
		}

		response, err := rasa.Query(c.Request.Text)
		if err != nil {
			c.Error(fmt.Sprintf("RASA ERROR %v", err))
			c.Next()
			return
		}

		if response.Intent.Name != "" {
			c.Flag("intent", strings.TrimSpace(response.Intent.Name))
		}

		// TODO Flag with entity results

		c.Flag("rasa", response)

		c.Next()
	}
}

// Query will make a query against the RASA api
func (l API) Query(queryString string) (*Response, error) {
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
		var rasaError interface{}

		err = json.Unmarshal(body, &rasaError)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("RASA Error %+v", rasaError)
	}

	rasaResponse := &Response{}

	err = json.Unmarshal(body, rasaResponse)
	if err != nil {
		return nil, err
	}

	return rasaResponse, nil

}
