package cpn

type MiddlewareFunction = func(c *Context) error

type NextFunction = func() error
