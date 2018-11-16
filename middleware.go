package gbl

type MiddlewareFunction = func(c *Context)

type NextFunction = func()
