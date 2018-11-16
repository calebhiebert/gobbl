package gbl

import (
	"testing"
)

func TestAddingMiddleware(t *testing.T) {
	bot := New()

	middleware := func(c *Context) {}

	bot.Use(middleware)

	if len(bot.middlewares) != 1 {
		t.Error("Bot is not storing middleware in it's internal slice")
	}
}
