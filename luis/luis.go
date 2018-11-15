/*
	luis.go

	This file contains methods for accessing the Microsoft LUIS API
	This middleware will take the c.Request.Text property and use it to query a LUIS endpoint
	The top scoring intent will be stored in the "intent" flag
	The entire result body will be stored in the "luis" flag
*/
package luis

import (
	"errors"
	"fmt"
	"net/url"

	".."
)

type LUIS struct {
	endpoint string

	queryParams map[string]string
}

func New(endpoint string) (*LUIS, error) {

	parsedEndpoint, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	config := LUIS{endpoint: endpoint, queryParams: make(map[string]string)}

	q := parsedEndpoint.Query()

	if _, exists := q["subscription-key"]; !exists {
		return nil, errors.New("Missing Subscription Key")
	}

	fmt.Println(parsedEndpoint.Host, parsedEndpoint.Path, parsedEndpoint.Query())

	return &config, nil
}

func LUISMiddleware(endpoint string) cpn.MiddlewareFunction {
	return func(c *cpn.Context) error {
		return nil
	}
}
