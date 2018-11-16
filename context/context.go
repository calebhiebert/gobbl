package bctx

import (
	"github.com/calebhiebert/gobbl"
)

type BotContext struct {
	Sequence int                  `json:"s"`
	Contexts map[string]BCContext `json:"c"`
}

// BCContext represents a single context item
type BCContext struct {
	Name            string `json:"n"`
	Source          string `json:"s"`
	Lifetime        int    `json:"l"`
	CurrentLifetime int    `json:"cl"`
}

func ClearAll(c *gbl.Context) {
	if c.HasFlag("_bctxDecoded") {
		botContext := c.GetFlag("_bctxDecoded").(*BotContext)
		botContext.Contexts = make(map[string]BCContext)
	} else {
		panic("Missing _bctxDecoded!")
	}
}

func Clear(c *gbl.Context, contextName string) {
	if c.HasFlag("_bctxDecoded") {
		botContext := c.GetFlag("_bctxDecoded").(*BotContext)
		delete(botContext.Contexts, contextName)
	} else {
		panic("Missing _bctxDecoded!")
	}
}

func Add(c *gbl.Context, context *BCContext) {
	if c.HasFlag("_bctxDecoded") {
		botContext := c.GetFlag("_bctxDecoded").(*BotContext)
		botContext.Contexts[context.Name] = *context
	} else {
		panic("Missing _bctxDecoded!")
	}
}
