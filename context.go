package gbl

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/matoous/go-nanoid"
)

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
	Identifier  string
	LogLevel    int
	abortErr    error
	logMutex    *sync.Mutex
	flagMutex   *sync.Mutex
}

type AbortFunction func(error)

// Transform Turns an input context struct into a full context
func (ic InputContext) Transform(bot *Bot) *Context {
	id, err := gonanoid.Nanoid(4)
	if err != nil {
		fmt.Println("ID GENERATION ERROR", err)
	}

	ctx := Context{
		RawRequest:  ic.RawRequest,
		Integration: ic.Integration,
		StartedAt:   time.Now().UnixNano(),
		Identifier:  id,
		R:           ic.Response,
		AutoRespond: true,
		LogLevel:    30,
		Flags:       make(map[string]interface{}),
		logMutex:    &sync.Mutex{},
		flagMutex:   &sync.Mutex{},
	}

	// Grab the log level from the environment
	if os.Getenv("LOG_LEVEL") != "" {
		ll, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
		if err != nil {
			ctx.Log(20, "The LOG_LEVEL environment variable must be a positive integer, not "+os.Getenv("LOG_LEVEL"), "LogLevelParser")
		} else {
			ctx.LogLevel = ll
		}
	}

	return &ctx
}

// Elapsed gets the number of milliseconds since the context was created
func (c Context) Elapsed() int64 {
	return (time.Now().UnixNano() - c.StartedAt) / 1000000
}

/*
	FLAG METHODS
*/

// Flag adds a flag to the context
func (c Context) Flag(key string, value interface{}) {
	c.flagMutex.Lock()
	c.Flags[key] = value
	c.flagMutex.Unlock()
}

// HasFlag returns true if a flag exists on the context
// false otherwise
func (c Context) HasFlag(key string) bool {
	c.flagMutex.Lock()
	_, exists := c.Flags[key]
	c.flagMutex.Unlock()

	return exists
}

func (c *Context) Abort(err error) {
	c.abortErr = err
}

// GetFlag will return the flag stored at key
func (c Context) GetFlag(key string) interface{} {
	defer c.flagMutex.Unlock()
	c.flagMutex.Lock()

	return c.Flags[key]
}

func (c Context) GetIntFlag(key string) int {
	defer c.flagMutex.Unlock()
	c.flagMutex.Lock()
	return c.Flags[key].(int)
}

func (c Context) GetInt8Flag(key string) int8 {
	defer c.flagMutex.Unlock()
	c.flagMutex.Lock()
	return c.Flags[key].(int8)
}

func (c Context) GetInt16Flag(key string) int16 {
	defer c.flagMutex.Unlock()
	c.flagMutex.Lock()
	return c.Flags[key].(int16)
}

func (c Context) GetInt32Flag(key string) int32 {
	defer c.flagMutex.Unlock()
	c.flagMutex.Lock()
	return c.Flags[key].(int32)
}

func (c Context) GetInt64Flag(key string) int64 {
	defer c.flagMutex.Unlock()
	c.flagMutex.Lock()
	return c.Flags[key].(int64)
}

func (c Context) GetStringFlag(key string) string {
	defer c.flagMutex.Unlock()
	c.flagMutex.Lock()
	return c.Flags[key].(string)
}

func (c Context) GetBoolFlag(key string) bool {
	defer c.flagMutex.Unlock()
	c.flagMutex.Lock()
	return c.Flags[key].(bool)
}

func (c Context) GetFloat64Flag(key string) float64 {
	defer c.flagMutex.Unlock()
	c.flagMutex.Lock()
	return c.Flags[key].(float64)
}

func (c Context) GetTimeFlag(key string) time.Time {
	defer c.flagMutex.Unlock()
	c.flagMutex.Lock()
	return c.Flags[key].(time.Time)
}

func (c Context) GetDurationFlag(key string) time.Duration {
	defer c.flagMutex.Unlock()
	c.flagMutex.Lock()
	return c.Flags[key].(time.Duration)
}

func (c Context) GetStringSliceFlag(key string) []string {
	defer c.flagMutex.Unlock()
	c.flagMutex.Lock()
	return c.Flags[key].([]string)
}

// ClearFlag will completely delete a flag from the context
func (c Context) ClearFlag(key ...string) {
	defer c.flagMutex.Unlock()
	c.flagMutex.Lock()

	for _, k := range key {
		delete(c.Flags, k)
	}
}
