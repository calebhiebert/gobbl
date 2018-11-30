// Package gbl contains some helpers for building a chatbot
package gbl

import (
	"fmt"
	"reflect"
	"runtime"
)

// Bot is a struct with a collection of middlewares
type Bot struct {
	middlewares  []MiddlewareFunction
	eventChan    chan Event
	eventHandler EventHandlerFunc
}

// DispatchFunction is a function used internally to run the middlewares
type DispatchFunction func(i int) error

// EventHandlerFunc is the function type that will handle bot events
type EventHandlerFunc func(event *Event)

// New creates a new bot instance
func New() *Bot {
	bot := Bot{
		middlewares: []MiddlewareFunction{},
		eventChan:   make(chan Event),
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

	if b.eventHandler != nil {
		b.eventChan <- Event{
			Type:    EVRequestStart,
			Context: preparedContext,
		}
	}

	err := b.exec(preparedContext)
	if err != nil {
		// Danger Danger
		return nil, err
	} else if preparedContext.abortErr != nil {
		preparedContext.Log(10, fmt.Sprintf("Request aborted %v", err), "Bot")
		return nil, preparedContext.abortErr
	}

	if b.eventHandler != nil {
		b.eventChan <- Event{
			Type:    EVRequestEnd,
			Context: preparedContext,
		}
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

		if b.eventHandler != nil {
			b.eventChan <- Event{
				Type: EVHandlerCall,
				HandlerCall: &HandlerCall{
					Handler:       runtime.FuncForPC(reflect.ValueOf(currentMiddleware).Pointer()).Name(),
					StackPosition: i,
				},
				Context: c,
			}
		}

		currentMiddleware(c)
		return nil
	}

	return dispatch(0)
}

// SetEventHandler sets the bot's event handler. This function will be called
// whenever a bot event is emitted
func (b *Bot) SetEventHandler(handler EventHandlerFunc) {
	if b.eventHandler != nil {
		panic("Event handler already set! An event handler can only be set once")
	}

	b.eventHandler = handler
	go handleEvents(b)
}

func handleEvents(bot *Bot) {
	for {
		event := <-bot.eventChan

		if bot.eventHandler != nil {
			bot.eventHandler(&event)
		}
	}
}
