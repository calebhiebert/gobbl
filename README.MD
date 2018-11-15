# GOBBL

GOBBL stands for **Go Bot Building Library**. It is a quick and easy way to get started with building a chatbot.

## Get Started

Here is the code for a simple echo bot

```go

func main() {

	// Create a new bot
	bot := gbl.New()

	// Use the request extraction middleware
	// to extract what the user has said from the context
	bot.Use(gbl.RequestExtractionMiddleware())

	// Add a simple middleware that will send an echo response
	bot.Use(func(c *gbl.Context) error {

		// When using the console integration, the context R (response) object
		// we need to cast it so we can use it's functions
		basicResponse := c.R.(*gbl.BasicResponse)

		// Add a text message to the output
		basicResponse.Text(fmt.Sprintf("Echo: %s", c.Request.Text))

		return nil
	})

	// Create a new console integration
	ci := gbl.ConsoleIntegration{}

	// Start listening to the console input
	ci.Listen(bot)
}


```

### Integrations

- Facebook Messenger
- LUIS

### Planned Integrations

- Telegram
- Apple Business Chat
- Chatfuel
- DialogFlow
- Rasa
