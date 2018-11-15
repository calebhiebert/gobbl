/*
	bot.go

	This file contains the main bot/middleware implementation
*/

// Package cpn contains some helpers for building a chatbot
package cpn

import (
	"fmt"
)

type Bot struct {
	middlewares []MiddlewareFunction
}

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
func (b *Bot) Execute(input *InputContext) (*[]Context, error) {
	preparedContext := input.Transform(b)

	err := b.exec(preparedContext)

	if err != nil {
		// Danger Danger
		return nil, err
	}

	if preparedContext.AutoRespond {
		_, err := preparedContext.Integration.Respond(preparedContext)
		if err != nil {
			fmt.Printf("Error while auto responding %+v", err)
		}
	}

	return &[]Context{*preparedContext}, nil
}

func (b *Bot) exec(c *Context) error {
	stackPosition := -1

	var dispatch DispatchFunction

	dispatch = func(i int) error {

		stackPosition = i

		if stackPosition == len(b.middlewares) {
			// We are at the bottom of the stack
			return nil
		}

		currentMiddleware := b.middlewares[i]

		c.Next = func() error {
			return dispatch(i + 1)
		}

		return currentMiddleware(c)
	}

	return dispatch(0)
}
