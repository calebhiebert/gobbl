package gbl

func SessionMiddleware() MiddlewareFunction {
	return func(c *Context) error {

		if c.sessionStore == nil {
			return c.Next()
		}

		session, err := c.sessionStore.Get(c.User.ID)
		if err != nil {
			return err
		}

		c.Session = *session

		err = c.Next()
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
