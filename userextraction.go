package gbl

func UserExtractionMiddleware() MiddlewareFunction {
	return func(c *Context) {
		user, err := c.Integration.User(c)
		if err != nil {
			c.Abort(err)
			return
		}

		c.User = user

		c.Next()
	}
}
