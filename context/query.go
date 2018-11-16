package bctx

type ContextQuery struct {
}

type ContextMatcherFunc func(ctx *BotContext) bool
