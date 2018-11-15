/*
	luis.go

	This file contains methods for accessing the Microsoft LUIS API
	This middleware will take the c.Request.Text property and use it to query a LUIS endpoint
	The top scoring intent will be stored in the "intent" flag
	The entire result body will be stored in the "luis" flag
*/
package luis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	".."
)

type LUIS struct {
	client   *http.Client
	endpoint string
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
		endpoint: fmt.Sprintf("https://%s%s%s", parsedEndpoint.Host, parsedEndpoint.Path, queryString),
		client: &http.Client{
			Timeout: 6 * time.Second,
		},
	}

	return &config, nil
}

// LUISMiddleware returns the LUIS middleware that will query luis with the Text property from the generic request
func LUISMiddleware(luis *LUIS) gbl.MiddlewareFunction {
	return func(c *gbl.Context) error {

		if c.Request.Text == "" {
			return c.Next()
		}

		response, err := luis.Query(c.Request.Text)
		if err != nil {
			fmt.Println(err)
			return c.Next()
		}

		fmt.Printf("LUIS %+v\n", response)

		if response.TopScoringIntent.Intent != "" {
			c.Flag("intent", response.TopScoringIntent.Intent)
		}

		c.Flag("luis", response)

		return c.Next()
	}
}

// Query will make a query against the LUIS api
func (l LUIS) Query(queryString string) (*LUISResponse, error) {
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

		return nil, errors.New(fmt.Sprintf("LUIS Error %+v", luisError))
	} else {
		luisResponse := &LUISResponse{}

		err = json.Unmarshal(body, luisResponse)
		if err != nil {
			return nil, err
		}

		return luisResponse, nil
	}
}
