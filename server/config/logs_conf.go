package config

import (
	"lib50cal/server/internal/logger"
	"log"
	"os"
)

// dummy struct to represent logger config
// there is no point for doing this, but we want some sort of integrity in config api
type LoggerConfig struct {
	print logger.Logger
	debug logger.Logger
}

// NewLoggerConfig returns a LoggerConfig which can be used to config logging service
// including debug and print for server package
func NewLoggerConfig() *LoggerConfig {
	// defualt debug and print logger, stdout for rainy days
	d := log.New(os.Stdout, "debug: ", log.LstdFlags|log.Lmsgprefix)
	p := log.New(os.Stdout, "print: ", log.LstdFlags|log.Lmsgprefix)

	logger.RegisterDebug(d)
	logger.RegisterPrint(p)
	return &LoggerConfig{debug: d, print: p}

}

// returns current debug logger
func (p *LoggerConfig) GetDebugLog() logger.Logger { return p.debug }

// returns current print logger
func (p *LoggerConfig) GetPrintLog() logger.Logger { return p.print }

// SetDebugLog sets new logger as debug for server package, if nil passed then logger pass
// all logs to /dev/null (io.Discard)
func (p *LoggerConfig) SetDebugLog(lg logger.Logger) {
	p.debug = lg
	logger.RegisterDebug(lg)
}

// SetPrintLog sets new logger as print for server package, if nil passed then logger pass
// all logs to /dev/null (io.Discard)
func (p *LoggerConfig) SetPrintLog(lg *log.Logger) {
	p.print = lg
	logger.RegisterPrint(lg)
}
