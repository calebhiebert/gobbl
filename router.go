/*
	router.go

	This file contains some basic router implimentations
*/

package gbl

// INTENT ROUTER
//
// This router will route to handlers based on the "intent" flag
type RIntentRouter struct {
	handlers map[string]MiddlewareFunction
}

// IntentRouter will create and return a new intent router
func IntentRouter() *RIntentRouter {
	r := RIntentRouter{
		handlers: make(map[string]MiddlewareFunction),
	}

	return &r
}

// Intent will add a new intent and handler pair to this router
func (r *RIntentRouter) Intent(intent string, handler MiddlewareFunction) {
	r.handlers[intent] = handler
}

// Middleware will return a middleware function that should be added to the bot
func (r *RIntentRouter) Middleware() MiddlewareFunction {
	return func(c *Context) error {

		if !c.HasFlag("intent") {
			return c.Next()
		}

		intent := c.GetFlag("intent").(string)

		handler, exists := r.handlers[intent]
		if !exists {
			return c.Next()
		} else {
			handler(c)
		}

		return nil
	}
}

// CUSTOM ROUTER
//
// This router will route to handlers based on a custom function
type RCustomRouter struct {
	pairs []CustomRouterPair
}

type CustomRouterPair struct {
	customFunc CustomRouterFunction
	handler    MiddlewareFunction
}
type CustomRouterFunction func(c *Context) bool

// CustomRouter will create and return a new custom router
func CustomRouter() *RCustomRouter {
	r := RCustomRouter{
		pairs: []CustomRouterPair{},
	}

	return &r
}

// Route will add a new intent and handler pair to this router
func (r *RCustomRouter) Route(customFunc CustomRouterFunction, handler MiddlewareFunction) {
	r.pairs = append(r.pairs, CustomRouterPair{
		customFunc: customFunc,
		handler:    handler,
	})
}

// Middleware will return a middleware function that should be added to the bot
func (r *RCustomRouter) Middleware() MiddlewareFunction {
	return func(c *Context) error {

		for _, routerPair := range r.pairs {
			if routerPair.customFunc(c) {
				return routerPair.handler(c)
			}
		}

		return c.Next()
	}
}
