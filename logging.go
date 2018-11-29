package gbl

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	aurora "github.com/logrusorgru/aurora"
)

// LogEntry represents a single log entry
type LogEntry struct {
	Elapsed int64
	Source  string
	Level   int
	Message string
}

// Log will log a statement to the console
func (c Context) Log(level int, msg, source string) {
	if level > c.LogLevel {
		return
	}

	if !c.HasFlag("__logs") {
		c.Flag("__logs", []LogEntry{})
	}

	elapsed := aurora.Green(fmt.Sprintf("+%dms", c.Elapsed())).Bold()
	formattedSource := aurora.Blue(source).String()
	formattedLevel := ""
	id := aurora.Magenta(c.Identifier).Bold()

	switch level {
	case 30:
		formattedLevel = aurora.Blue("INFO").String()
	case 40:
		formattedLevel = aurora.Cyan("DEBUG").String()
	case 50:
		formattedLevel = aurora.Gray("TRACE").String()
	case 20:
		formattedLevel = aurora.Brown("WARN").String()
	case 10:
		formattedLevel = aurora.Red("ERROR").String()
	default:
		formattedLevel = aurora.Gray("CLVL " + strconv.Itoa(level)).String()
	}

	fmt.Printf("[%s | %s | %s] %s %s\n", elapsed, formattedSource, id, formattedLevel, msg)

	logEntry := LogEntry{
		Elapsed: c.Elapsed(),
		Source:  source,
		Level:   level,
		Message: msg,
	}

	c.logMutex.Lock()
	logArr := c.GetFlag("__logs").([]LogEntry)
	logArr = append(logArr, logEntry)
	c.Flag("__logs", logArr)
	c.logMutex.Unlock()
}

// Info will log a statement at the INFO level
func (c Context) Info(msg string) {
	c.Log(30, msg, GetCallingFunction())
}

// Debug will log a statement at the DEBUG level
func (c Context) Debug(msg string) {
	c.Log(40, msg, GetCallingFunction())
}

// Warn will log a statement at the WARN level
func (c Context) Warn(msg string) {
	c.Log(20, msg, GetCallingFunction())
}

// Error will log a statement at the ERROR level
func (c Context) Error(msg string) {
	c.Log(10, msg, GetCallingFunction())
}

// Trace will log a statement at the TRACE level
func (c Context) Trace(msg string) {
	c.Log(50, msg, GetCallingFunction())
}

// Infof will log an info statement with foratted args
func (c Context) Infof(format string, args ...interface{}) {
	c.Log(30, fmt.Sprintf(format, args...), GetCallingFunction())
}

// Debugf will log an debug statement with foratted args
func (c Context) Debugf(format string, args ...interface{}) {
	c.Log(40, fmt.Sprintf(format, args...), GetCallingFunction())
}

// Tracef will log an debug statement with foratted args
func (c Context) Tracef(format string, args ...interface{}) {
	c.Log(50, fmt.Sprintf(format, args...), GetCallingFunction())
}

// Warnf will log an debug statement with foratted args
func (c Context) Warnf(format string, args ...interface{}) {
	c.Log(20, fmt.Sprintf(format, args...), GetCallingFunction())
}

// Errorf will log an debug statement with foratted args
func (c Context) Errorf(format string, args ...interface{}) {
	c.Log(10, fmt.Sprintf(format, args...), GetCallingFunction())
}

// GetCallingFunction will return the name of the function that called
// the function that calls this function
func GetCallingFunction() string {
	fpcs := make([]uintptr, 1)

	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return "n/a"
	}

	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "n/a"
	}

	nameParts := strings.Split(fun.Name(), ".")

	return nameParts[len(nameParts)-1]
}
