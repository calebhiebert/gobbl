package gbl

import "time"

type InputContext struct {
	RawRequest  interface{}
	Integration Integration
	Response    interface{}
}

type Context struct {
	RawRequest  interface{}
	User        User
	Integration Integration
	AutoRespond bool
	R           interface{}
	Request     GenericRequest
	StartedAt   int64
	Flags       map[string]interface{}
	Next        NextFunction
	abortErr    error
}

type AbortFunction func(error)

// Turns an input context struct into a full context
func (ic InputContext) Transform(bot *Bot) *Context {
	ctx := Context{
		RawRequest:  ic.RawRequest,
		Integration: ic.Integration,
		StartedAt:   time.Now().Unix(),
		R:           ic.Response,
		AutoRespond: true,
		Flags:       make(map[string]interface{}),
	}

	return &ctx
}

// Gets the number of milliseconds since the context was created
func (c Context) Elapsed() int64 {
	return time.Now().Unix() - c.StartedAt
}

/*
	FLAG METHODS
*/
func (c Context) Flag(key string, value interface{}) {
	c.Flags[key] = value
}

func (c Context) HasFlag(key string) bool {
	_, exists := c.Flags[key]

	return exists
}

func (c *Context) Abort(err error) {
	c.abortErr = err
}

func (c Context) GetFlag(key string) interface{} {
	return c.Flags[key]
}

func (c Context) GetIntFlag(key string) int {
	return c.Flags[key].(int)
}

func (c Context) GetInt8Flag(key string) int8 {
	return c.Flags[key].(int8)
}

func (c Context) GetInt16Flag(key string) int16 {
	return c.Flags[key].(int16)
}

func (c Context) GetInt32Flag(key string) int32 {
	return c.Flags[key].(int32)
}

func (c Context) GetInt64Flag(key string) int64 {
	return c.Flags[key].(int64)
}

func (c Context) GetStringFlag(key string) string {
	return c.Flags[key].(string)
}

func (c Context) GetBoolFlag(key string) bool {
	return c.Flags[key].(bool)
}

func (c Context) GetFloat64Flag(key string) float64 {
	return c.Flags[key].(float64)
}

func (c Context) GetTimeFlag(key string) time.Time {
	return c.Flags[key].(time.Time)
}

func (c Context) GetDurationFlag(key string) time.Duration {
	return c.Flags[key].(time.Duration)
}

func (c Context) GetStringSliceFlag(key string) []string {
	return c.Flags[key].([]string)
}

func (c Context) ClearFlag(key string) {
	delete(c.Flags, key)
}
