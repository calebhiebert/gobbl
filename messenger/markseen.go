/*
	markseen.go

	This file contains a middleware function that will send the mark seen user action to facebook
	If no user ID is present, the middleware will just call c.Next()

	Note: This middleware should only be added after the user extraction middleware is added
*/
package fb

import (
	"fmt"

	"github.com/calebhiebert/gobbl"
)

// MarkSeenMiddleware will mark the current chat as seen by the bot
func MarkSeenMiddleware() gbl.MiddlewareFunction {
	return func(c *gbl.Context) {
		if c.User.ID != "" {
			go func() {
				_, err := c.Integration.(*MessengerIntegration).API.SenderAction(&User{
					ID: c.User.ID,
				}, SenderActionMarkSeen)
				if err != nil {
					c.Error(fmt.Sprintf("Error when setting as seen %v", err))
				}
			}()
		}

		c.Next()
	}
}
