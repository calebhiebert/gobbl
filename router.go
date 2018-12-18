/*
	router.go

	This file contains some basic router implimentations
*/

package gbl

// RIntentRouter will route to handlers based on the "intent" flag
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

// Intents will add multiple intents that point to a single handler
func (r *RIntentRouter) Intents(intents []string, handler MiddlewareFunction) {
	for _, intent := range intents {
		r.Intent(intent, handler)
	}
}

// Middleware will return a middleware function that should be added to the bot
func (r *RIntentRouter) Middleware() MiddlewareFunction {
	return func(c *Context) {

		if !c.HasFlag("intent") {
			c.Next()
			return
		}

		intent := c.GetFlag("intent").(string)

		handler, exists := r.handlers[intent]
		if !exists {
			c.Next()
		} else {
			handler(c)
		}
	}
}

// RCustomRouter will route to handlers based on a custom function
type RCustomRouter struct {
	pairs []customRouterPair
}

type customRouterPair struct {
	customFunc CustomRouterFunction
	handler    MiddlewareFunction
}

// CustomRouterFunction is a function that will return true if the router
// should route to the current handler, false otherwise
type CustomRouterFunction func(c *Context) bool

// CustomRouter will create and return a new custom router
func CustomRouter() *RCustomRouter {
	r := RCustomRouter{
		pairs: []customRouterPair{},
	}

	return &r
}

// Route will add a new intent and handler pair to this router
func (r *RCustomRouter) Route(customFunc CustomRouterFunction, handler MiddlewareFunction) {
	r.pairs = append(r.pairs, customRouterPair{
		customFunc: customFunc,
		handler:    handler,
	})
}

// Middleware will return a middleware function that should be added to the bot
func (r *RCustomRouter) Middleware() MiddlewareFunction {
	return func(c *Context) {

		for _, routerPair := range r.pairs {
			if routerPair.customFunc(c) {
				routerPair.handler(c)
				return
			}
		}

		c.Next()
	}
}

// RTextRouter will route to handlers based on the request text
type RTextRouter struct {
	handlers map[string]MiddlewareFunction
}

// TextRouter will create and return a new intent router
func TextRouter() *RTextRouter {
	r := RTextRouter{
		handlers: make(map[string]MiddlewareFunction),
	}

	return &r
}

// Text will add a new text and handler pair to this router
func (r *RTextRouter) Text(text string, handler MiddlewareFunction) {
	r.handlers[text] = handler
}

// Middleware will return a middleware function that should be added to the bot
func (r *RTextRouter) Middleware() MiddlewareFunction {
	return func(c *Context) {
		if c.Request.Text == "" {
			c.Next()
			return
		}

		text := c.Request.Text

		handler, exists := r.handlers[text]
		if !exists {
			c.Next()
		} else {
			handler(c)
		}
	}
}
