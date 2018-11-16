package sess

import (
	"strings"

	"github.com/calebhiebert/gobbl"
)

func Middleware(store SessionStore) gbl.MiddlewareFunction {
	return func(c *gbl.Context) error {

		session, err := store.Get(c.User.ID)
		if err != nil {
			if err == ErrSessionNonexistant {
				session = make(map[string]interface{})
			} else {
				return err
			}
		}

		populateSessionFlags(c, session)

		err = c.Next()
		if err != nil {
			return err
		}

		sessionToSave := readSessionFlags(c)

		err = store.Update(c.User.ID, &sessionToSave)
		if err != nil {
			return err
		}

		return nil
	}
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
			flags[k] = v
		}
	}

	return flags
}
