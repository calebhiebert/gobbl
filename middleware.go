package cpn

type MiddlewareFunction = func(c *Context, next NextFunction) error

type NextFunction = func() error
