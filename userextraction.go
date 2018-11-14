package cpn

func UserExtractionMiddleware() MiddlewareFunction {
	return func(c *Context, next NextFunction) error {

		user, err := c.Integration.User(c)
		if err != nil {
			return err
		}

		c.User = user

		return next()
	}
}
