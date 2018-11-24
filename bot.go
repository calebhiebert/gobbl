// Package gbl contains some helpers for building a chatbot
package gbl

// Bot is a struct with a collection of middlewares
type Bot struct {
	middlewares []MiddlewareFunction
}

// DispatchFunction is a function used internally to run the middlewares
type DispatchFunction func(i int) error

// New creates a new bot instance
func New() *Bot {
	bot := Bot{
		middlewares: []MiddlewareFunction{},
	}

	return &bot
}

// Use will add a new middleware to the bot.
// Middlewares will be executed in the order that they were added
func (b *Bot) Use(f MiddlewareFunction) {
	b.middlewares = append(b.middlewares, f)
}

// Execute will take the InputContext and generate a full context from it.
// It will take this full context
func (b *Bot) Execute(input *InputContext) (*Context, error) {
	preparedContext := input.Transform(b)

	err := b.exec(preparedContext)
	if err != nil {
		// Danger Danger
		return nil, err
	} else if preparedContext.abortErr != nil {
		return nil, preparedContext.abortErr
	}

	return preparedContext, nil
}

func (b *Bot) exec(c *Context) error {
	stackPosition := -1

	var dispatch DispatchFunction

	dispatch = func(i int) error {

		if c.abortErr != nil {
			return c.abortErr
		}

		stackPosition = i

		if stackPosition == len(b.middlewares) {
			// We are at the bottom of the stack
			return nil
		}

		currentMiddleware := b.middlewares[i]

		c.Next = func() {
			dispatch(i + 1)
		}

		currentMiddleware(c)
		return nil
	}

	return dispatch(0)
}
