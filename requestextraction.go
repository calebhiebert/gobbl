package cpn

func RequestExtractionMiddleware() MiddlewareFunction {
	return func(c *Context) error {

		req, err := c.Integration.GenericRequest(c)
		if err != nil {
			return err
		}

		c.Request = *req

		return c.Next()
	}
}
