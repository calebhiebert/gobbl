package bctx

import (
	"fmt"

	"github.com/calebhiebert/gobbl"
)

// Middleware will generate the context middleware, this will take care
// of setting and managing user sessions
func Middleware() gbl.MiddlewareFunction {
	return func(c *gbl.Context) {

		// Check to see if the context is set in the session
		if !c.HasFlag("sess:_bctx") {

			// Create a new context for the session
			blankContext := BotContext{
				Contexts: make(map[string]BCContext),
			}

			// Encode our new blank context
			encodedContext, err := encodeContext(&blankContext)
			if err != nil {
				c.Error(fmt.Sprintf("Context Error %v", err))
			}

			// Flag the session with it
			c.Flag("sess:_bctx", encodedContext)
		}

		// Decode the context from the session
		decodedContext, err := decodeContext(c.GetStringFlag("sess:_bctx"))
		if err != nil {
			c.Error(fmt.Sprintf("Context Decode Error %v", err))
			c.Next()
			return
		}

		decodedContext.Sequence++

		// Create a new slice to store all the contexts that are still alive
		liveContexts := map[string]BCContext{}

		// Increment context current life time and track living ones
		for name, contextEntry := range decodedContext.Contexts {

			if contextEntry.CurrentLifetime > 0 && contextEntry.BirthSequence != decodedContext.Sequence-1 {
				contextEntry.CurrentLifetime--
			}

			if contextEntry.CurrentLifetime > 0 || contextEntry.CurrentLifetime == -1 {
				liveContexts[name] = contextEntry
			}
		}

		decodedContext.Contexts = liveContexts

		c.Flag(flagKeyName, &decodedContext)

		// Complete the bot runthrough
		c.Next()

		// Encode any changes to the context and set it on the session
		if !c.HasFlag(flagKeyName) {
			c.Error(fmt.Sprintf("_bctxDecoded flag was removed! Context will be unusable"))
			return
		}

		updatedContext := c.GetFlag(flagKeyName).(*BotContext)

		encodedUpdatedContext, err := encodeContext(updatedContext)
		if err != nil {
			c.Error(fmt.Sprintf("Error encoding updated context %v", err))
		}

		c.Flag("sess:_bctx", encodedUpdatedContext)
	}
}
