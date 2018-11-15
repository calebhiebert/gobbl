package gbl

type MiddlewareFunction = func(c *Context) error

type NextFunction = func() error
