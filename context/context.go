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
	BirthSequence   int    `json:"bs"`
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

func Add(c *gbl.Context, name string, lifetime int) {
	AddSourced(c, name, lifetime, "none")
}

func AddSourced(c *gbl.Context, name string, lifetime int, source string) {
	if c.HasFlag("_bctxDecoded") {
		botContext := c.GetFlag("_bctxDecoded").(*BotContext)
		botContext.Contexts[name] = BCContext{
			Name:            name,
			BirthSequence:   botContext.Sequence,
			CurrentLifetime: lifetime,
			Source:          source,
		}
	} else {
		panic("Missing _bctxDecoded!")
	}
}
