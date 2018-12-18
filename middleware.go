package gbl

// MiddlewareFunction type is a gobbl handler
type MiddlewareFunction = func(c *Context)

// NextFunction is what c.Next() is
type NextFunction = func()
