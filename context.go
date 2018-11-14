package cap

import "time"

type InputContext struct {
	RawRequest  interface{}
	Integration Integration
	Response    interface{}
}

type Context struct {
	RawRequest   *interface{}
	User         *User
	Integration  Integration
	AutoRespond  bool
	R            interface{}
	Request      GenericRequest
	StartedAt    int64
	Session      map[string]interface{}
	sessionStore SessionStore
}

// Turns an input context struct into a full context
func (ic InputContext) Transform(bot *Bot) *Context {
	ctx := Context{
		RawRequest:  &ic.RawRequest,
		Integration: ic.Integration,
		StartedAt:   time.Now().Unix(),
		R:           ic.Response,
		AutoRespond: true,
	}

	return &ctx
}

func (c Context) SaveSession() error {
	return c.sessionStore.Update(c.User.ID, &c.Session)
}

// Gets the number of milliseconds since the context was created
func (c Context) Elapsed() int64 {
	return time.Now().Unix() - c.StartedAt
}
