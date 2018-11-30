package gbl

import "errors"

// ErrMultipleNext is returned when the c.Next() function has been called multiple times
var ErrMultipleNext = errors.New("Next() called mulitple times")
