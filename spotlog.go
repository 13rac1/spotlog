package spotlog

import (
	"context"

	"github.com/sirupsen/logrus"
)

type contextKey string

const loggerKey contextKey = "spotlogger"

// Get returns the logger in the context or creates one.
func Get(ctx context.Context) (context.Context, *Logger) {
	logger, ok := ctx.Value(loggerKey).(*Logger)

	if ok {
		return ctx, logger
	}

	log := logrus.New()
	log.Level = logrus.DebugLevel
	logger = &Logger{
		Logger: log,
		level:  logrus.ErrorLevel,
	}
	ctx = context.WithValue(ctx, loggerKey, logger)
	return ctx, logger
}

type storedEntry struct {
	level logrus.Level
	args  []interface{}
}

type Logger struct {
	*logrus.Logger
	level   logrus.Level
	entries []storedEntry
}

func (logger *Logger) Log(level logrus.Level, args ...interface{}) {
	if level <= logger.level {
		for _, entry := range logger.entries {
			logger.Logger.Log(entry.level, entry.args...)
		}
		// Clear the list of entries output.
		logger.entries = nil

		logger.Logger.Log(level, args...)
	} else {
		logger.entries = append(logger.entries, storedEntry{level, args})
	}
}

func (logger *Logger) Trace(args ...interface{}) {
	logger.Log(logrus.TraceLevel, args...)
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.Log(logrus.DebugLevel, args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.Log(logrus.InfoLevel, args...)
}

// func (logger *Logger) Print(args ...interface{}) {
// 	entry := logger.newEntry()
// 	entry.Print(args...)
// 	logger.releaseEntry(entry)
// }

func (logger *Logger) Warn(args ...interface{}) {
	logger.Log(logrus.WarnLevel, args...)
}

func (logger *Logger) Warning(args ...interface{}) {
	logger.Warn(args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.Log(logrus.ErrorLevel, args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.Log(logrus.FatalLevel, args...)
	logger.Exit(1)
}

func (logger *Logger) Panic(args ...interface{}) {
	logger.Log(logrus.PanicLevel, args...)
}
