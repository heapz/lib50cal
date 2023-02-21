package server

import (
	"io"
	"log"
	"os"
	"sync/atomic"
)

// dummy struct to represent logger config
// there is no point for doing this, but we want some sort of integrity in config api
type loggerConfig struct{}

// crazy ass user may call NewLoggerConfig multiple times, we need law and order here
var atdebug = atomic.Pointer[log.Logger]{}
var atprint = atomic.Pointer[log.Logger]{}

// NewLoggerConfig returns a loggerConfig which can be used to config logging service
// including debug and print for server package
func NewLoggerConfig() *loggerConfig {

	// defualt debug and print logger, stdout for rainy days
	atdebug.Store(log.New(os.Stdout, "debug: ", log.LstdFlags|log.Lmsgprefix))
	atprint.Store(log.New(os.Stdout, "print: ", log.LstdFlags|log.Lmsgprefix))
	return &loggerConfig{}

}

// returns current debug logger
func (p *loggerConfig) GetDebugLog() *log.Logger { return atdebug.Load() }

// returns current print logger
func (p *loggerConfig) GetPrintLog() *log.Logger { return atprint.Load() }

// SetDebugLog sets new logger as debug for server package, if nil passed then logger pass
// all logs to /dev/null (io.Discard)
func (p *loggerConfig) SetDebugLog(logger *log.Logger) {
	if logger == nil {
		atdebug.Load().SetOutput(io.Discard)
		return
	}

	atdebug.Swap(logger)
}

// SetPrintLog sets new logger as print for server package, if nil passed then logger pass
// all logs to /dev/null (io.Discard)
func (p *loggerConfig) SetPrintLog(logger *log.Logger) {
	if logger == nil {
		atprint.Load().SetOutput(io.Discard)
		return
	}

	atprint.Swap(logger)
}
