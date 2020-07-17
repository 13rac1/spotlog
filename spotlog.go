package spotlog

import (
	"context"

	"github.com/sirupsen/logrus"
)

type contextKey string
type printType string

const (
	loggerKey  contextKey = "spotlogger"
	printLog   printType  = "log"
	printLogf  printType  = "logf"
	printLogln printType  = "logln"
)

// storedEntry contains the arguments of stored Log Entry.
type storedEntry struct {
	method printType
	level  logrus.Level
	format string
	args   []interface{}
}

// Get returns the logger in the context or creates one.
func Get(ctx context.Context) (context.Context, *Logger) {
	logger, ok := ctx.Value(loggerKey).(*Logger)

	if ok {
		return ctx, logger
	}

	logrusLogger := logrus.New()
	// The logrus logger is set to TraceLevel to print everything.
	logrusLogger.Level = logrus.TraceLevel

	// TODO: Is there a way to allow a global Logger instance and store the
	// entries in the Context?
	logger = &Logger{
		Logger:      logrusLogger,
		entries:     []storedEntry{},
		minLogLevel: logrus.ErrorLevel,
	}
	ctx = context.WithValue(ctx, loggerKey, logger)

	return ctx, logger
}
