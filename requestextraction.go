package main

func RequestExtractionMiddleware() MiddlewareFunction {
	return func(c *Context, next NextFunction) error {

		req, err := c.Integration.GenericRequest(c)
		if err != nil {
			return err
		}

		c.Request = *req

		return next()
	}
}
