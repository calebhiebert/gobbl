package bctx

import (
	"github.com/calebhiebert/gobbl"
)

var flagKeyName = "_bctxDecoded"

// BotContext is the var that gets set on the session
type BotContext struct {
	Sequence int                  `json:"s"`
	Contexts map[string]BCContext `json:"c"`
}

// BCContext represents a single context item
type BCContext struct {
	Name            string            `json:"n"`
	Source          string            `json:"s"`
	BirthSequence   int               `json:"bs"`
	CurrentLifetime int               `json:"cl"`
	Data            map[string]string `json:"d"`
}

// ClearAll will clear all contexts for the current session
func ClearAll(c *gbl.Context) {
	if c.HasFlag(flagKeyName) {
		botContext := c.GetFlag(flagKeyName).(*BotContext)
		botContext.Contexts = make(map[string]BCContext)
	} else {
		panic("Missing _bctxDecoded!")
	}
}

// Clear will clear the context with the given name from the session (if it exists)
func Clear(c *gbl.Context, contextName string) {
	if c.HasFlag(flagKeyName) {
		botContext := c.GetFlag(flagKeyName).(*BotContext)
		delete(botContext.Contexts, contextName)
	} else {
		panic("Missing _bctxDecoded!")
	}
}

// Add will add a new context to the session
func Add(c *gbl.Context, name string, lifetime int) {
	AddSourced(c, name, lifetime, "none")
}

// AddSourced will add a new context to the session with an optional source param
// this is mostly helpful for debugging, since you can see where the context was added
func AddSourced(c *gbl.Context, name string, lifetime int, source string) {
	if c.HasFlag(flagKeyName) {
		botContext := c.GetFlag(flagKeyName).(*BotContext)
		botContext.Contexts[name] = BCContext{
			Name:            name,
			BirthSequence:   botContext.Sequence,
			CurrentLifetime: lifetime,
			Source:          source,
			Data:            make(map[string]string),
		}
	} else {
		panic("Missing _bctxDecoded!")
	}
}

// Get will return the context data param stored at the context key value pair location
func Get(c *gbl.Context, contextName, dataParam string) string {
	if c.HasFlag(flagKeyName) {
		botContext := c.GetFlag(flagKeyName).(*BotContext)

		ctx, exists := botContext.Contexts[contextName]
		if !exists {
			return ""
		}

		return ctx.Data[dataParam]
	}

	return ""
}

// Set will set a key value pair on the context
func Set(c *gbl.Context, contextName, dataParam, dataValue string) {
	if c.HasFlag(flagKeyName) {
		botContext := c.GetFlag(flagKeyName).(*BotContext)

		ctx, exists := botContext.Contexts[contextName]
		if exists {
			ctx.Data[dataParam] = dataValue
		}
	}
}
