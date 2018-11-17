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

func MarkSeenMiddleware() gbl.MiddlewareFunction {
	return func(c *gbl.Context) {
		if c.User.ID != "" {
			go func() {
				_, err := c.Integration.(*MessengerIntegration).API.SenderAction(&User{
					ID: c.User.ID,
				}, SenderActionMarkSeen)
				if err != nil {
					fmt.Println("SET TYPING ERR", err)
				}
			}()
		}

		c.Next()
	}
}
