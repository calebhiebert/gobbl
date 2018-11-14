package main

type Bot struct {
	middlewares []MiddlewareFunction
}

type DispatchFunction func(i int) error

func (b Bot) Use(f MiddlewareFunction) {
	b.middlewares = append(b.middlewares, f)
}

func (b Bot) Execute(input *InputContext) (*[]Context, error) {
	preparedContext := prepareContext(input)

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

// Turns a context input into a prepared context
func prepareContext(input *InputContext) *Context {

	ctx := Context{}

	return &ctx
}
