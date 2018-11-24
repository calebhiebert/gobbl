package gbl

import (
	"fmt"
)

// ResponderMiddleware will use the current integration to respond
// to the incoming message
func ResponderMiddleware() MiddlewareFunction {
	return func(c *Context) {
		// We don't want to respond right away, so we can just wait
		// for the request to finish
		c.Next()

		// Only respond if AutoRespond is true
		if c.AutoRespond {
			_, err := c.Integration.Respond(c)
			if err != nil {
				c.Error(fmt.Sprintf("Error during responding process! %v", err))
			}
		}
	}
}
