package bctx

import (
	"github.com/calebhiebert/gobbl"
)

type RContextIntentRouter struct {
	handlers  map[string][]ContextRouterQuery
	fallbacks []ContextRouterQuery
}

type ContextRouterQuery struct {
	matcherFunc ContextMatcherFunc
	intentOnly  bool
	noContext   bool
	any         C
	all         C
	handler     gbl.MiddlewareFunction
}

type I []string
type C []string

func ContextIntentRouter() *RContextIntentRouter {
	return &RContextIntentRouter{
		handlers:  make(map[string][]ContextRouterQuery),
		fallbacks: make([]ContextRouterQuery, 0),
	}
}

func (cr *RContextIntentRouter) IntentsOnly(intents I, handler gbl.MiddlewareFunction) {
	for _, intent := range intents {
		addQueryForIntent(cr, intent, &ContextRouterQuery{
			intentOnly: true,
			handler:    handler,
		})
	}
}

func (cr *RContextIntentRouter) NoContext(intents I, handler gbl.MiddlewareFunction) {
	for _, intent := range intents {
		addQueryForIntent(cr, intent, &ContextRouterQuery{
			noContext: true,
			handler:   handler,
		})
	}
}

func (cr *RContextIntentRouter) FallbackAny(contexts C, handler gbl.MiddlewareFunction) {
	cr.fallbacks = append(cr.fallbacks, ContextRouterQuery{
		handler: handler,
		any:     contexts,
	})
}

func (cr *RContextIntentRouter) FallbackAll(contexts C, handler gbl.MiddlewareFunction) {
	cr.fallbacks = append(cr.fallbacks, ContextRouterQuery{
		handler: handler,
		all:     contexts,
	})
}

func (cr *RContextIntentRouter) Any(intents I, contexts C, handler gbl.MiddlewareFunction) {
	for _, intent := range intents {
		addQueryForIntent(cr, intent, &ContextRouterQuery{
			any:     contexts,
			handler: handler,
		})
	}
}

func (cr *RContextIntentRouter) All(intents I, contexts C, handler gbl.MiddlewareFunction) {
	for _, intent := range intents {
		addQueryForIntent(cr, intent, &ContextRouterQuery{
			all:     contexts,
			handler: handler,
		})
	}
}

func addQueryForIntent(cr *RContextIntentRouter, intent string, query *ContextRouterQuery) {
	if queries, exists := cr.handlers[intent]; !exists {
		cr.handlers[intent] = []ContextRouterQuery{*query}
	} else {
		cr.handlers[intent] = append(queries, *query)
	}
}

func hasAllContexts(botContext *BotContext, contexts []string) bool {
	for _, ctx := range contexts {
		if _, exists := botContext.Contexts[ctx]; !exists {
			return false
		}
	}

	return true
}

func hasAnyContexts(botContext *BotContext, contexts []string) bool {
	for _, ctx := range contexts {
		if _, exists := botContext.Contexts[ctx]; exists {
			return true
		}
	}

	return false
}

// Middleware generates the middleware required to use this router
func (cr *RContextIntentRouter) Middleware() gbl.MiddlewareFunction {
	return func(c *gbl.Context) {
		if c.HasFlag("_bctxDecoded") {
			botContext := c.GetFlag("_bctxDecoded").(*BotContext)

			if c.HasFlag("intent") {
				intent := c.GetStringFlag("intent")

				if intentCollection, exists := cr.handlers[intent]; exists {
					for _, query := range intentCollection {
						if query.intentOnly {
							query.handler(c)
							return
						} else if query.all != nil && hasAllContexts(botContext, query.all) {
							query.handler(c)
							return
						} else if query.any != nil && hasAnyContexts(botContext, query.any) {
							query.handler(c)
							return
						}
					}
				}
			}

			for _, fallback := range cr.fallbacks {
				if fallback.all != nil && hasAllContexts(botContext, fallback.all) {
					fallback.handler(c)
					return
				} else if fallback.any != nil && hasAnyContexts(botContext, fallback.any) {
					fallback.handler(c)
					return
				}
			}

			if c.HasFlag("intent") {
				intent := c.GetStringFlag("intent")

				if intentCollection, exists := cr.handlers[intent]; exists && len(botContext.Contexts) == 0 {
					for _, query := range intentCollection {
						if query.noContext {
							query.handler(c)
							return
						}
					}
				}
			}
		} else {
			panic("Missing _bctxDecoded!")
		}

		c.Next()
	}
}
