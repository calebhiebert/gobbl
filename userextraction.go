package gbl

// User is an object that can be extracted from each request
// not all fields will be populated. At the very least, the ID field
// should be populated
type User struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
}

// UserExtractionMiddleware will extract a user out of the incoming request
func UserExtractionMiddleware() MiddlewareFunction {
	return func(c *Context) {
		user, err := c.Integration.User(c)
		if err != nil {
			c.Errorf("Unable to extract user %v", err)
		} else {
			c.User = user
		}

		c.Next()
	}
}
