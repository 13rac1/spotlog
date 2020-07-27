package spotlog

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type Entry struct {
	*logrus.Entry
	Logger *SpotLogger
}

func NewEntry(logger *SpotLogger) *Entry {
	return &Entry{
		Entry:  logrus.NewEntry(logger.Logger),
		Logger: logger,
	}
}

// Add an error as single field (using the key defined in ErrorKey) to the Entry.
func (e *Entry) WithError(err error) *Entry {
	return e.WithField(logrus.ErrorKey, err)
}

// Add a context to the Entry.
func (e *Entry) WithContext(ctx context.Context) *Entry {
	dataCopy := make(logrus.Fields, len(e.Data))
	for k, v := range e.Data {
		dataCopy[k] = v
	}
	return &Entry{Logger: e.Logger, Entry: e.Entry.WithContext(ctx)}
}

// Add a single field to the Entry.
func (e *Entry) WithField(key string, value interface{}) *Entry {
	return e.WithFields(logrus.Fields{key: value})
}

// Add a map of fields to the Entry.
func (e *Entry) WithFields(fields logrus.Fields) *Entry {
	return &Entry{Logger: e.Logger, Entry: e.Entry.WithFields(fields)}
}

// Overrides the time of the Entry.
func (e *Entry) WithTime(t time.Time) *Entry {
	return &Entry{Logger: e.Logger, Entry: e.Entry.WithTime(t)}
}

// Original comment from logrus: This function is not declared with a pointer
// value because otherwise race conditions will occur when using multiple
// goroutines.
func (e Entry) log(method printType, level logrus.Level, format string, args ...interface{}) {
	e.Logger.entriesLock.Lock()
	defer e.Logger.entriesLock.Unlock()

	if e.Logger.alwaysLog(level) {
		// Found an important log, print the stored log entries.
		for _, entry := range e.Logger.entries {
			switch method {
			case printLog:
				e.Entry.Log(entry.level, entry.args...)
			case printLogf:
				e.Entry.Logf(entry.level, entry.format, entry.args...)
			case printLogln:
				e.Entry.Logln(entry.level, entry.args...)
			}
		}

		// Clear the list of output entries.
		e.Logger.entries = nil
		// Then print the actual "important" log.
		e.Entry.Log(level, args...)
	} else {
		e.Logger.entries = append(e.Logger.entries, storedEntry{method, level, format, args})
	}
}

func (e *Entry) Log(level logrus.Level, args ...interface{}) {
	e.log(printLog, level, "", args...)
}

func (e *Entry) Trace(args ...interface{}) {
	e.Log(logrus.TraceLevel, args...)
}

func (e *Entry) Debug(args ...interface{}) {
	e.Log(logrus.DebugLevel, args...)
}

func (e *Entry) Print(args ...interface{}) {
	e.Info(args...)
}

func (e *Entry) Info(args ...interface{}) {
	e.Log(logrus.InfoLevel, args...)
}

func (e *Entry) Warn(args ...interface{}) {
	e.Log(logrus.WarnLevel, args...)
}

func (e *Entry) Warning(args ...interface{}) {
	e.Warn(args...)
}

func (e *Entry) Error(args ...interface{}) {
	e.Log(logrus.ErrorLevel, args...)
}

func (e *Entry) Fatal(args ...interface{}) {
	e.Log(logrus.FatalLevel, args...)
	e.Logger.Exit(1)
}

func (e *Entry) Panic(args ...interface{}) {
	e.Log(logrus.PanicLevel, args...)
	panic(fmt.Sprint(args...))
}

// Entry Printf family functions

func (e *Entry) Logf(level logrus.Level, format string, args ...interface{}) {
	e.log(printLogf, level, "", args...)
}

func (e *Entry) Tracef(format string, args ...interface{}) {
	e.Logf(logrus.TraceLevel, format, args...)
}

func (e *Entry) Debugf(format string, args ...interface{}) {
	e.Logf(logrus.DebugLevel, format, args...)
}

func (e *Entry) Infof(format string, args ...interface{}) {
	e.Logf(logrus.InfoLevel, format, args...)
}

func (e *Entry) Printf(format string, args ...interface{}) {
	e.Infof(format, args...)
}

func (e *Entry) Warnf(format string, args ...interface{}) {
	e.Logf(logrus.WarnLevel, format, args...)
}

func (e *Entry) Warningf(format string, args ...interface{}) {
	e.Warnf(format, args...)
}

func (e *Entry) Errorf(format string, args ...interface{}) {
	e.Logf(logrus.ErrorLevel, format, args...)
}

func (e *Entry) Fatalf(format string, args ...interface{}) {
	e.Logf(logrus.FatalLevel, format, args...)
	e.Logger.Exit(1)
}

func (e *Entry) Panicf(format string, args ...interface{}) {
	e.Logf(logrus.PanicLevel, format, args...)
}

// Entry Println family functions

func (e *Entry) Logln(level logrus.Level, args ...interface{}) {
	e.log(printLogln, level, "", args...)
}

func (e *Entry) Traceln(args ...interface{}) {
	e.Logln(logrus.TraceLevel, args...)
}

func (e *Entry) Debugln(args ...interface{}) {
	e.Logln(logrus.DebugLevel, args...)
}

func (e *Entry) Infoln(args ...interface{}) {
	e.Logln(logrus.InfoLevel, args...)
}

func (e *Entry) Println(args ...interface{}) {
	e.Infoln(args...)
}

func (e *Entry) Warnln(args ...interface{}) {
	e.Logln(logrus.WarnLevel, args...)
}

func (e *Entry) Warningln(args ...interface{}) {
	e.Warnln(args...)
}

func (e *Entry) Errorln(args ...interface{}) {
	e.Logln(logrus.ErrorLevel, args...)
}

func (e *Entry) Fatalln(args ...interface{}) {
	e.Logln(logrus.FatalLevel, args...)
	e.Logger.Exit(1)
}

func (e *Entry) Panicln(args ...interface{}) {
	e.Logln(logrus.PanicLevel, args...)
}
