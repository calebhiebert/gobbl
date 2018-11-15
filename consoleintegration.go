package cpn

import (
	"bufio"
	"fmt"
	"os"
)

type ConsoleIntegration struct {
}

func (ci *ConsoleIntegration) GenericRequest(c *Context) (*GenericRequest, error) {
	genericRequest := GenericRequest{
		Text: c.RawRequest.(string),
	}

	return &genericRequest, nil
}

func (ci *ConsoleIntegration) User(c *Context) (User, error) {
	user := User{
		ID:        "consoleid",
		FirstName: "John",
		LastName:  "Smith",
		Email:     "john.smith@dummymail.com",
	}

	return user, nil
}

func (ci *ConsoleIntegration) Respond(c *Context) (*interface{}, error) {
	for _, msg := range c.R.(*BasicResponse).messages {
		fmt.Printf("[bot] %s\n", msg)
	}

	return nil, nil
}

func (ci *ConsoleIntegration) Listen(bot *Bot) {
	reader := bufio.NewReader(os.Stdin)

	var input = ""
	var err error

	for input != "exit()" {
		fmt.Print("> ")
		input, err = reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		inputCtx := InputContext{
			RawRequest:  input,
			Integration: ci,
			Response:    &BasicResponse{},
		}

		_, err = bot.Execute(&inputCtx)
		if err != nil {
			panic(err)
		}
	}

}
