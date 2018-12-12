package bctx

import (
	"github.com/calebhiebert/gobbl"
)

// RContextIntentRouter is a router that can be used to route requests
// based on the current user context and intent
type RContextIntentRouter struct {
	handlers  map[string][]ContextRouterQuery
	fallbacks []ContextRouterQuery
}

// ContextRouterQuery represents a set of parameters used to decide
// to fire this handler or not
type ContextRouterQuery struct {
	intentOnly bool
	noContext  bool
	any        C
	all        C
	handler    gbl.MiddlewareFunction
}

// I is a type alias for a list of intents
type I []string

// C is a type alias for a list of contexts
type C []string

// ContextIntentRouter will create a new context intent router
func ContextIntentRouter() *RContextIntentRouter {
	return &RContextIntentRouter{
		handlers:  make(map[string][]ContextRouterQuery),
		fallbacks: make([]ContextRouterQuery, 0),
	}
}

// NoContext will match a route if the intent matches and the user currently has no
// contexts at all
func (cr *RContextIntentRouter) NoContext(intents I, handler gbl.MiddlewareFunction) {
	for _, intent := range intents {
		addQueryForIntent(cr, intent, &ContextRouterQuery{
			noContext: true,
			handler:   handler,
		})
	}
}

// FallbackAny will match if the user has any of the provided contexts, regardless of intent
func (cr *RContextIntentRouter) FallbackAny(contexts C, handler gbl.MiddlewareFunction) {
	cr.fallbacks = append(cr.fallbacks, ContextRouterQuery{
		handler: handler,
		any:     contexts,
	})
}

// FallbackAll will match if the user has all of the provided contexts, regardless of intent
func (cr *RContextIntentRouter) FallbackAll(contexts C, handler gbl.MiddlewareFunction) {
	cr.fallbacks = append(cr.fallbacks, ContextRouterQuery{
		handler: handler,
		all:     contexts,
	})
}

// Any will match if the intent matches, and any of the supplied contexts are present
func (cr *RContextIntentRouter) Any(intents I, contexts C, handler gbl.MiddlewareFunction) {
	for _, intent := range intents {
		addQueryForIntent(cr, intent, &ContextRouterQuery{
			any:     contexts,
			handler: handler,
		})
	}
}

// All will match if the intent matches, and all of the supplied contexts are present
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

		// Only run checks if the context object exists on the request
		if c.HasFlag("_bctxDecoded") {
			botContext := c.GetFlag("_bctxDecoded").(*BotContext)

			// First check if an intent is present
			if c.HasFlag("intent") {
				intent := c.GetStringFlag("intent")

				// If an intent is present, check to see if any handlers match
				if intentCollection, exists := cr.handlers[intent]; exists {
					for _, query := range intentCollection {
						if query.all != nil && hasAllContexts(botContext, query.all) {
							query.handler(c)
							return
						} else if query.any != nil && hasAnyContexts(botContext, query.any) {
							query.handler(c)
							return
						}
					}
				}
			}

			// Process fallbacks, only after intents have been checked
			for _, fallback := range cr.fallbacks {
				if fallback.all != nil && hasAllContexts(botContext, fallback.all) {
					fallback.handler(c)
					return
				} else if fallback.any != nil && hasAnyContexts(botContext, fallback.any) {
					fallback.handler(c)
					return
				}
			}

			// Check no context intents, only after in-context intents, and fallbacks have been processed
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
