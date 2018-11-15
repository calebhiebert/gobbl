package gbl

type Integration interface {

	// Extract a "generic request" from the full request, the generic request is in a format that everything in the bot can understand
	GenericRequest(c *Context) (GenericRequest, error)

	// Extract a user object from the request
	User(c *Context) (User, error)

	// Uses the information in the context to respond to the request
	Respond(c *Context) (*interface{}, error)
}
