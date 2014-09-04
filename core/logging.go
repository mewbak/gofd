package core

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// logger, the default singleton logger
var logger = Logger{prefix: "", logging_level: 3, out: os.Stdout}
var DefLogger = &logger

const LOG_DEBUG = 0
const LOG_INFO = 1
const LOG_ERROR = 2
const LOG_NONE = 3

// Logger, struct containing information for constructing log messages
type Logger struct {
	mu            sync.Mutex // ensures atomic I/O; protects other fields
	prefix        string     // prefix to write at beginning of each line
	logging_level int
	out           io.Writer
}

// GetLogger returns the default singleton logger
func GetLogger() *Logger {
	return DefLogger
}

// output prints with exklusive access to the output stream
func (this *Logger) output(msg string) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.out.Write(getByteArrFromString(msg))
}

// SetLoggingLevel sets the logging level
func (this *Logger) SetLoggingLevel(level int) {
	this.logging_level = level
}

// GetLoggingLevel returns the current logging level
func (this *Logger) GetLoggingLevel() int {
	return this.logging_level
}

// getByteArrFromString returns string contents as byte array,
// helper for stdout, which only accepts byte array
func getByteArrFromString(s string) []byte {
	var buf []byte
	return append(buf, s...)
}

// I prints message on info log level
func (this *Logger) I(msg string) {
	if LOG_INFO >= this.logging_level {
		this.output(fmt.Sprintf("INFO: %v", msg))
	}
}

// Iln prints message with newline on info log level
func (this *Logger) Iln(msg string) {
	if LOG_INFO >= this.logging_level {
		this.output(fmt.Sprintf("INFO: %v\n", msg))
	}
}

// If prints message via Printf on info log level
func (this *Logger) If(format string, v ...interface{}) {
	if LOG_INFO >= this.logging_level {
		s := fmt.Sprintf(format, v...)
		this.output(fmt.Sprintf("INFO: %v\n", s))
	}
}

// D prints message on debug log level
func (this *Logger) D(msg string) {
	if LOG_DEBUG >= this.logging_level {
		this.output(fmt.Sprintf("DEBUG: %v", msg))
	}
}

// Dln prints message with newline on debug log level
func (this *Logger) Dln(msg string) {
	if LOG_DEBUG >= this.logging_level {
		this.output(fmt.Sprintf("DEBUG: %v\n", msg))
	}
}

// Df prints message via Printf on debug log level
func (this *Logger) Df(format string, v ...interface{}) {
	if LOG_DEBUG >= this.logging_level {
		s := fmt.Sprintf(format, v...)
		this.output(fmt.Sprintf("DEBUG: %v\n", s))
	}
}

// E prints message on error log level
func (this *Logger) E(msg string) {
	if LOG_ERROR >= this.logging_level {
		this.output(fmt.Sprintf("ERROR: %v", msg))
	}
}

// Eln prints message with newline on error log level
func (this *Logger) Eln(msg string) {
	if LOG_ERROR >= this.logging_level {
		this.output(fmt.Sprintf("ERROR: %v\n", msg))
	}
}

// Ef prints message via Printf on error log level
func (this *Logger) Ef(format string, v ...interface{}) {
	if LOG_ERROR >= this.logging_level {
		s := fmt.Sprintf(format, v...)
		this.output(fmt.Sprintf("ERROR: %v\n", s))
	}
}

// P prints message regardless of log level
func (this *Logger) P(msg string) {
	this.output(msg)
}

// Pln prints message with newline regardless of log level
func (this *Logger) Pln(msg string) {
	this.output(msg)
}

// Pf prints message via Printf regardless of log level
func (this *Logger) Pf(format string, v ...interface{}) {
	this.output(fmt.Sprintf(format+"\n", v...))
}

// DoInfo checks whether log level is LOG_INFO
func (this *Logger) DoInfo() bool {
	return this.GetLoggingLevel() == LOG_INFO
}

// DoDebug checks whether log level is LOG_DEBUG
func (this *Logger) DoDebug() bool {
	return this.GetLoggingLevel() == LOG_DEBUG
}

type ILogger interface {
	I(msg string)
	D(msg string)
	E(msg string)
	Iln(msg string)
	Dln(msg string)
	Eln(msg string)
	If(format string, v ...interface{})
	Df(format string, v ...interface{})
	Ef(format string, v ...interface{})
	P(msg string)
	Pln(msg string)
	Pf(format string, v ...interface{})
	SetLoggingLevel(val int)
	GetLoggingLevel() int
}

/*
// Examples:
func main() {
	logger.SetLoggingLevel(LOG_INFO)
	logger.I("blub")
	logger.Iln("bla")
	logger.If("bla: %v %v\n","blub","blubber")
	logger.D("blub")
	logger.Dln("bla")
	logger.Df("bla: %v %v\n","blub","blubber")
	logger.E("blub")
	logger.Eln("bla")
	logger.Ef("bla: %v %v\n","blub","blubber")
	logger.P("blub")
	logger.Pln("bla")
	logger.Pf("bla: %v %v\n","blub","blubber")
	logger.I("blub")
}
*/
