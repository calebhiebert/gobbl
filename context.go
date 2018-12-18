package gbl

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/matoous/go-nanoid"
)

// InputContext represents the arguments available when executing
// an incoming request on the bot
type InputContext struct {
	RawRequest  interface{}
	Integration Integration
	Response    interface{}
}

// Context is the GOBBL context object
type Context struct {
	RawRequest  interface{}            `json:"raw"`
	User        User                   `json:"user"`
	Integration Integration            `json:"-"`
	AutoRespond bool                   `json:"autoRespond"`
	R           interface{}            `json:"res"`
	Request     GenericRequest         `json:"req"`
	StartedAt   int64                  `json:"startedAt"`
	Flags       map[string]interface{} `json:"-"`
	Next        NextFunction           `json:"-"`
	Identifier  string                 `json:"id"`
	LogLevel    int                    `json:"logLevel"`
	abortErr    error
	logMutex    *sync.Mutex
	flagMutex   *sync.Mutex
	bot         *Bot
}

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
		bot:         bot,
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

	if c.bot.eventHandler != nil {
		c.bot.eventChan <- Event{
			Type: EVFlagSet,
			FlagSet: &FlagSet{
				Flag:  key,
				Value: fmt.Sprintf("%+v", value),
			},
			Context: &c,
		}
	}
}

// HasFlag returns true if a flag exists on the context
// false otherwise
func (c Context) HasFlag(key string) bool {
	c.flagMutex.Lock()
	_, exists := c.Flags[key]
	c.flagMutex.Unlock()

	if c.bot.eventHandler != nil {
		c.bot.eventChan <- Event{
			Type: EVFlagAccess,
			FlagAccess: &FlagAccess{
				Flag:             key,
				IsExistenceCheck: true,
			},
			Context: &c,
		}
	}

	return exists
}

// GetFlag will return the flag stored at key
func (c Context) GetFlag(key string) interface{} {
	defer c.flagMutex.Unlock()

	if c.bot.eventHandler != nil {
		c.bot.eventChan <- Event{
			Type: EVFlagAccess,
			FlagAccess: &FlagAccess{
				Flag:             key,
				IsExistenceCheck: false,
			},
			Context: &c,
		}
	}

	c.flagMutex.Lock()

	return c.Flags[key]
}

func (c Context) GetIntFlag(key string) int {
	return c.GetFlag(key).(int)
}

func (c Context) GetInt8Flag(key string) int8 {
	return c.GetFlag(key).(int8)
}

func (c Context) GetInt16Flag(key string) int16 {
	return c.GetFlag(key).(int16)
}

func (c Context) GetInt32Flag(key string) int32 {
	return c.GetFlag(key).(int32)
}

func (c Context) GetInt64Flag(key string) int64 {
	return c.GetFlag(key).(int64)
}

func (c Context) GetStringFlag(key string) string {
	return c.GetFlag(key).(string)
}

func (c Context) GetBoolFlag(key string) bool {
	return c.GetFlag(key).(bool)
}

func (c Context) GetFloat64Flag(key string) float64 {
	return c.GetFlag(key).(float64)
}

func (c Context) GetTimeFlag(key string) time.Time {
	return c.GetFlag(key).(time.Time)
}

func (c Context) GetDurationFlag(key string) time.Duration {
	return c.GetFlag(key).(time.Duration)
}

func (c Context) GetStringSliceFlag(key string) []string {
	slice := c.GetFlag(key)

	switch slice.(type) {
	case []string:
		return slice.([]string)
	case []interface{}:
		newSlice := []string{}
		for _, v := range slice.([]interface{}) {
			newSlice = append(newSlice, v.(string))
		}
		return newSlice
	}

	return slice.([]string)
}

// ClearFlag will completely delete a flag from the context
func (c Context) ClearFlag(key ...string) {
	defer c.flagMutex.Unlock()

	if c.bot.eventHandler != nil {
		c.bot.eventChan <- Event{
			Type: EVFlagClear,
			FlagClear: &FlagClear{
				Flags: key,
			},
			Context: &c,
		}
	}

	c.flagMutex.Lock()

	for _, k := range key {
		delete(c.Flags, k)
	}
}
