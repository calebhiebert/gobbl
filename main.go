package main

func main() {
	bot := New()

	bot.Use(UserExtractionMiddleware())
	bot.Use(RequestExtractionMiddleware())
	bot.Use(SessionMiddleware())

	bot.Use(func(c *Context, next NextFunction) error {
		c.R.(*BasicResponse).Text("This is a text")

		return next()
	})

	ci := ConsoleIntegration{}

	ci.Listen(bot)
}
