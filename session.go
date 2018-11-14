package cap

func SessionMiddleware() MiddlewareFunction {
	return func(c *Context, next NextFunction) error {

		if c.sessionStore == nil {
			return next()
		}

		session, err := c.sessionStore.Get(c.User.ID)
		if err != nil {
			return err
		}

		c.Session = *session

		err = next()
		if err != nil {
			return err
		}

		err = c.SaveSession()
		if err != nil {
			return err
		}

		return nil
	}
}
