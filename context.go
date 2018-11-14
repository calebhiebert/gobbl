package main

type InputContext struct {
	RawRequest *interface{}
}

type Context struct {
	RawRequest   *interface{}
	User         User
	Integration  Integration
	AutoRespond  bool
	R            interface{}
	Request      GenericRequest
	StartedAt    int64
	Session      map[string]interface{}
	sessionStore SessionStore
}

func (c Context) SaveSession() error {
	return c.sessionStore.Update(c.User.ID, &c.Session)
}
