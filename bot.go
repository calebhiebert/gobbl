package main

type Bot struct {
	middlewares []MiddlewareFunction
}

type DispatchFunction func(i int) error

func New() *Bot {
	bot := Bot{
		middlewares: make([]MiddlewareFunction, 0),
	}

	return &bot
}

func (b Bot) Use(f MiddlewareFunction) {
	b.middlewares = append(b.middlewares, f)
}

func (b Bot) Execute(input *InputContext) (*[]Context, error) {
	preparedContext := input.Transform(&b)

	err := b.exec(preparedContext)

	if err != nil {
		// Danger Danger
		return nil, err
	}

	return &[]Context{*preparedContext}, nil
}

func (b Bot) exec(c *Context) error {
	stackPosition := -1

	var dispatch DispatchFunction

	dispatch = func(i int) error {

		stackPosition = i

		if stackPosition == len(b.middlewares) {
			// We are at the bottom of the stack
			return nil
		}

		currentMiddleware := b.middlewares[i]

		return currentMiddleware(c, func() error {
			return dispatch(i + 1)
		})
	}

	return dispatch(0)
}
