package gbl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ConsoleIntegration provides the gobbl integration interface
type ConsoleIntegration struct {
}

// GenericRequest will grab the user's console input
func (ci *ConsoleIntegration) GenericRequest(c *Context) (GenericRequest, error) {
	genericRequest := GenericRequest{
		Text: c.RawRequest.(string),
	}

	return genericRequest, nil
}

// User will provide a preset user
func (ci *ConsoleIntegration) User(c *Context) (User, error) {
	user := User{
		ID:        "consoleid",
		FirstName: "John",
		LastName:  "Smith",
		Email:     "john.smith@dummymail.com",
	}

	return user, nil
}

// Respond will print a message to the console
func (ci *ConsoleIntegration) Respond(c *Context) (*interface{}, error) {
	if c.R != nil {
		fmt.Printf("[bot] %s\n", c.R)
	}

	return nil, nil
}

// Listen for incoming console messages
func (ci *ConsoleIntegration) Listen(bot *Bot) {
	reader := bufio.NewReader(os.Stdin)

	var input = ""
	var err error

	for input != "exit()" {
		func() {
			fmt.Print("> ")
			input, err = reader.ReadString('\n')
			if err != nil {
				panic(err)
			}

			input = strings.TrimSpace(input)

			inputCtx := InputContext{
				RawRequest:  strings.TrimSpace(input),
				Integration: ci,
				Response:    nil,
			}

			_, err = bot.Execute(&inputCtx)
			if err != nil {
				fmt.Printf("[ERROR] %s\n", err.Error())
			}
		}()
	}
}
