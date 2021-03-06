<p align="center">
  <img width="282" height="250" src="https://user-images.githubusercontent.com/16126072/49407072-054f5e80-f71d-11e8-82a5-2f867d031cfd.png">
</p>

# GOBBL

GOBBL stands for **Go Bot Building Library**. It is a quick and easy way to get started with building a chatbot.

## Get Started

First, install the library

```bash
go get -u github.com/calebhiebert/gobbl
```

Here is the code for a simple echo bot

```go

func main() {

	// Create a new bot
	bot := gbl.New()
	
	// Use this middleware to make sure the bot responds to requests 
	gobblr.Use(gbl.ResponderMiddleware())

	// Use the request extraction middleware
	// to extract what the user has said from the context
	bot.Use(gbl.RequestExtractionMiddleware())


	// Add a simple middleware that will send an echo response
	bot.Use(func(c *gbl.Context) {

		// When using the console integration, the context R (response) object
		// we need to cast it so we can use it's functions
		basicResponse := c.R.(*gbl.BasicResponse)

		// Add a text message to the output
		basicResponse.Text(fmt.Sprintf("Echo: %s", c.Request.Text))
	})

	// Create a new console integration
	ci := gbl.ConsoleIntegration{}

	// Start listening to the console input
	ci.Listen(bot)
}
```

## More Information

For more documentation check out the [Wiki](https://github.com/calebhiebert/gobbl/wiki)

### Integrations

- Facebook Messenger
- LUIS
- Rasa
- Telegram

### Planned Integrations

- Apple Business Chat
- Chatfuel
- DialogFlow
