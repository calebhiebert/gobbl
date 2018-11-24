package gbl

import (
	"fmt"
	"runtime"
	"strings"
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
	abortErr    error
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
		StartedAt:   time.Now().Unix(),
		Identifier:  id,
		R:           ic.Response,
		AutoRespond: true,
		Flags:       make(map[string]interface{}),
	}

	return &ctx
}

// Elapsed gets the number of milliseconds since the context was created
func (c Context) Elapsed() int64 {
	return time.Now().Unix() - c.StartedAt
}

/*
	FLAG METHODS
*/

// Flag adds a flag to the context
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

func (c Context) ClearFlag(key ...string) {
	for _, k := range key {
		delete(c.Flags, k)
	}
}

// Log will log a statement to the console
func (c Context) Log(level, msg string) {
	fmt.Printf("[+%dms - %s - %s] %s %s", c.Elapsed(), GetCallingFunction(), c.Identifier, level, msg)
}

// Info will log a statement at the INFO level
func (c Context) Info(msg string) {
	c.Log("INFO", msg)
}

// Debug will log a statement at the DEBUG level
func (c Context) Debug(msg string) {
	c.Log("DEBUG", msg)
}

// Warn will log a statement at the WARN level
func (c Context) Warn(msg string) {
	c.Log("WARN", msg)
}

// Error will log a statement at the ERROR level
func (c Context) Error(msg string) {
	c.Log("ERROR", msg)
}

// Trace will log a statement at the TRACE level
func (c Context) Trace(msg string) {
	c.Log("TRACE", msg)
}

// GetCallingFunction will return the name of the function that called
// the function that calls this function
func GetCallingFunction() string {
	fpcs := make([]uintptr, 1)

	n := runtime.Callers(4, fpcs)
	if n == 0 {
		return "n/a"
	}

	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "n/a"
	}

	nameParts := strings.Split(fun.Name(), ".")

	return nameParts[len(nameParts)-1]
}
