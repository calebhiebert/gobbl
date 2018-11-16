package gbl

func RequestExtractionMiddleware() MiddlewareFunction {
	return func(c *Context) {

		req, err := c.Integration.GenericRequest(c)
		if err != nil {
			c.Abort(err)
			return
		}

		c.Request = req

		c.Next()
	}
}
