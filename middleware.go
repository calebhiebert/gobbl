package main

type MiddlewareFunction = func(c *Context, next NextFunction) error

type NextFunction = func() error
