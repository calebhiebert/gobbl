package bctx

import (
	"fmt"

	"github.com/calebhiebert/gobbl"
)

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
				fmt.Println("Context Error", err)
			}

			// Flag the session with it
			c.Flag("sess:_bctx", encodedContext)
		}

		// Decode the context from the session
		decodedContext, err := decodeContext(c.GetStringFlag("sess:_bctx"))
		if err != nil {
			fmt.Println("Context Decode Error", err)
			c.Next()
			return
		}

		decodedContext.Sequence++

		// Create a new slice to store all the contexts that are still alive
		liveContexts := map[string]BCContext{}

		// Increment context current life time and track living ones
		for name, contextEntry := range decodedContext.Contexts {

			if contextEntry.CurrentLifetime > 0 {
				contextEntry.CurrentLifetime--
			}

			if contextEntry.CurrentLifetime > 0 {
				liveContexts[name] = contextEntry
			}
		}

		decodedContext.Contexts = liveContexts

		c.Flag("_bctxDecoded", &decodedContext)

		// Complete the bot runthrough
		c.Next()

		// Encode any changes to the context and set it on the session
		if !c.HasFlag("_bctxDecoded") {
			fmt.Println("_bctxDecoded flag was removed! Context will be unusable")
			return
		}

		updatedContext := c.GetFlag("_bctxDecoded").(*BotContext)

		encodedUpdatedContext, err := encodeContext(updatedContext)
		if err != nil {
			fmt.Println("Error encoding updated context")
		}

		c.Flag("sess:_bctx", encodedUpdatedContext)
	}
}
