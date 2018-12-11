package sess

import (
	"errors"
	"strings"

	"github.com/calebhiebert/gobbl"
)

var ErrSessionNonexistant = errors.New("Session did not exist")

func Middleware(store SessionStore) gbl.MiddlewareFunction {
	return func(c *gbl.Context) {

		session, err := store.Get(c.User.ID)
		if err != nil {
			if err == ErrSessionNonexistant {
				session = make(map[string]interface{})
			}
		}

		populateSessionFlags(c, session)

		// Wait for the request to finish
		c.Next()

		sessionToSave := readSessionFlags(c)

		err = store.Update(c.User.ID, &sessionToSave)
		if err != nil {
			c.Abort(err)
		}
	}
}

// ClearSession will clear all session variables
func ClearSession(c *gbl.Context) {
	flags := []string{}

	for k := range c.Flags {
		if strings.HasPrefix(k, "sess:") {
			flags = append(flags, k)
		}
	}

	c.ClearFlag(flags...)
}

func populateSessionFlags(c *gbl.Context, data map[string]interface{}) {
	for k, v := range data {
		c.Flag("sess:"+k, v)
	}
}

func readSessionFlags(c *gbl.Context) map[string]interface{} {

	flags := map[string]interface{}{}

	for k, v := range c.Flags {
		if strings.HasPrefix(k, "sess:") {
			flags[k[5:]] = v
		}
	}

	return flags
}
