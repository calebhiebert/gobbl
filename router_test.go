/*
	router_test.go

	This file contains the unit tests for the routers
*/

package gbl

import "testing"

func TestIntentRouterAddHandler(t *testing.T) {
	ir := IntentRouter()

	presetHandler := func(c *Context) error {
		return nil
	}

	ir.Intent("test-intent", presetHandler)

	_, present := ir.handlers["test-intent"]
	if !present {
		t.Error("Handler was not present on router!")
	}
}
