package cpn

func UserExtractionMiddleware() MiddlewareFunction {
	return func(c *Context) error {
		user, err := c.Integration.User(c)
		if err != nil {
			return err
		}

		c.User = user

		return c.Next()
	}
}
