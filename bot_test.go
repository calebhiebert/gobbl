package cpn

import (
	"testing"
)

func TestAddingMiddleware(t *testing.T) {
	bot := New()

	middleware := func(c *Context, next NextFunction) error {
		return nil
	}

	bot.Use(middleware)

	if len(bot.middlewares) != 1 {
		t.Error("Bot is not storing middleware in it's internal slice")
	}
}
