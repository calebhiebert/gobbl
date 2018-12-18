package gbl

// GenericRequest represents data that should be extractable from
// every incoming request
type GenericRequest struct {
	Text string
}

// RequestExtractionMiddleware will use the current integration to extract
// generic parameters from the incoming request
func RequestExtractionMiddleware() MiddlewareFunction {
	return func(c *Context) {

		req, err := c.Integration.GenericRequest(c)
		if err != nil {
			c.Errorf("Unable to extract generic request %v", err)
		} else {
			c.Request = req
		}

		c.Next()
	}
}
