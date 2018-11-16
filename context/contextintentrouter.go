package bctx

import (
	"github.com/calebhiebert/gobbl"
)

type RContextIntentRouter struct {
	handlers map[string][]ContextRouterPair
}

type ContextRouterPair struct {
	matcherFunc ContextMatcherFunc
	handler     gbl.MiddlewareFunction
}

func (cr *RContextIntentRouter) ICtx(intent string, match ContextMatcherFunc, handler gbl.MiddlewareFunction) {
	if intentPairs, exists := cr.handlers[intent]; !exists {

		pair := ContextRouterPair{
			matcherFunc: match,
			handler:     handler,
		}

		cr.handlers[intent] = []ContextRouterPair{pair}
	} else {
		cr.handlers[intent] = append(intentPairs, ContextRouterPair{
			matcherFunc: match,
			handler:     handler,
		})
	}
}

func (cr *RContextIntentRouter) Middleware() gbl.MiddlewareFunction {
	return func(c *gbl.Context) {
		if !c.HasFlag("intent") {
			c.Next()
			return
		}

		if c.HasFlag("_bctxDecoded") {
			botContext := c.GetFlag("_bctxDecoded").(*BotContext)

			if intentCollection, exists := cr.handlers[c.GetStringFlag("intent")]; exists {
				for _, routerPair := range intentCollection {
					if routerPair.matcherFunc(botContext) {
						routerPair.handler(c)
						return
					}
				}
			}
		} else {
			panic("Missing _bctxDecoded!")
		}

		c.Next()
	}
}
