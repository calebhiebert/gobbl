package bctx

import (
	"github.com/calebhiebert/gobbl"
)

type RContextRouter struct {
	pairs []contextRouterPair
}

type contextRouterPair struct {
	matcherFunc ContextMatcherFunc
	handler     gbl.MiddlewareFunction
}

func (cr *RContextRouter) Ctx(match ContextMatcherFunc, handler gbl.MiddlewareFunction) {
	cr.pairs = append(cr.pairs, contextRouterPair{
		matcherFunc: match,
		handler:     handler,
	})
}

func (cr *RContextRouter) Middleware() gbl.MiddlewareFunction {
	return func(c *gbl.Context) {

		if c.HasFlag("_bctxDecoded") {
			botContext := c.GetFlag("_bctxDecoded").(*BotContext)

			for _, routerPair := range cr.pairs {
				if routerPair.matcherFunc(botContext) {
					routerPair.handler(c)
					return
				}
			}

		} else {
			panic("Missing _bctxDecoded!")
		}

		c.Next()
	}
}
