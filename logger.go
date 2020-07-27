package spotlog

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// New creates a configured SpotLogger.
func New() *SpotLogger {
	logrusLogger := logrus.StandardLogger()
	// The logrus logger is set to TraceLevel to print everything.
	logrusLogger.Level = logrus.TraceLevel

	return &SpotLogger{
		Logger:      logrusLogger,
		entries:     []storedEntry{},
		minLogLevel: logrus.ErrorLevel,
	}
}

// SpotLogger wraps logrus.Logger to add log storage.
type SpotLogger struct {
	*logrus.Logger
	// minLogLevel is the minimum log level to output.
	minLogLevel logrus.Level

	entries     []storedEntry
	entriesLock sync.Mutex
}

func (l *SpotLogger) alwaysLog(level logrus.Level) bool {
	// Levels have lower values the higher their priority is.
	return level <= l.minLogLevel
}

func (l *SpotLogger) newEntry() *Entry {
	// TODO: Use Pool
	// entry, ok := l.entryPool.Get().(*Entry)
	// if ok {
	// 	return entry
	// }
	return NewEntry(l)
}

func (l *SpotLogger) releaseEntry(entry *Entry) {
	//entry.Data = map[string]interface{}{}
	//l.entryPool.Put(entry)
}

// WithField allocates a new entry and adds a field to it.
// Debug, Print, Info, Warn, Error, Fatal or Panic must be then applied to
// this new returned entry.
// If you want multiple fields, use `WithFields`.
func (l *SpotLogger) WithField(key string, value interface{}) *Entry {
	entry := l.newEntry()
	defer l.releaseEntry(entry)
	return entry.WithField(key, value)
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (l *SpotLogger) WithFields(fields logrus.Fields) *Entry {
	entry := l.newEntry()
	defer l.releaseEntry(entry)
	return entry.WithFields(fields)
}

// Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
func (l *SpotLogger) WithError(err error) *Entry {
	entry := l.newEntry()
	defer l.releaseEntry(entry)
	return entry.WithError(err)
}

// Add a context to the log entry.
func (l *SpotLogger) WithContext(ctx context.Context) *Entry {
	entry := l.newEntry()
	defer l.releaseEntry(entry)
	return entry.WithContext(ctx)
}

// Overrides the time of the log entry.
func (l *SpotLogger) WithTime(t time.Time) *Entry {
	entry := l.newEntry()
	defer l.releaseEntry(entry)
	return entry.WithTime(t)
}

func (l *SpotLogger) log(method printType, level logrus.Level, format string, args ...interface{}) {
	l.entriesLock.Lock()
	defer l.entriesLock.Unlock()

	if l.alwaysLog(level) {
		// Found an important log, print the stored log entries.
		for _, entry := range l.entries {
			switch method {
			case printLog:
				l.Logger.Log(entry.level, entry.args...)
			case printLogf:
				l.Logger.Logf(entry.level, entry.format, entry.args...)
			case printLogln:
				l.Logger.Logln(entry.level, entry.args...)
			}
		}

		// Clear the list of output entries.
		l.entries = nil
		// Then print the actual "important" log entry.
		l.Logger.Log(level, args...)
	} else {
		l.entries = append(l.entries, storedEntry{method, level, format, args})
	}
}

func (l *SpotLogger) Logf(level logrus.Level, format string, args ...interface{}) {
	l.log(printLogf, level, format, args...)
}

func (l *SpotLogger) Tracef(format string, args ...interface{}) {
	l.Logf(logrus.TraceLevel, format, args...)
}

func (l *SpotLogger) Debugf(format string, args ...interface{}) {
	l.Logf(logrus.DebugLevel, format, args...)
}

func (l *SpotLogger) Infof(format string, args ...interface{}) {
	l.Logf(logrus.InfoLevel, format, args...)
}

func (l *SpotLogger) Printf(format string, args ...interface{}) {
	l.Logf(logrus.InfoLevel, format, args...)
}

func (l *SpotLogger) Warnf(format string, args ...interface{}) {
	l.Logf(logrus.WarnLevel, format, args...)
}

func (l *SpotLogger) Warningf(format string, args ...interface{}) {
	l.Warnf(format, args...)
}

func (l *SpotLogger) Errorf(format string, args ...interface{}) {
	l.Logf(logrus.ErrorLevel, format, args...)
}

func (l *SpotLogger) Fatalf(format string, args ...interface{}) {
	l.Logf(logrus.FatalLevel, format, args...)
	l.Exit(1)
}

func (l *SpotLogger) Panicf(format string, args ...interface{}) {
	l.Logf(logrus.PanicLevel, format, args...)
}

func (l *SpotLogger) Log(level logrus.Level, args ...interface{}) {
	l.log(printLog, level, "", args...)
}

func (l *SpotLogger) Trace(args ...interface{}) {
	l.Log(logrus.TraceLevel, args...)
}

func (l *SpotLogger) Debug(args ...interface{}) {
	l.Log(logrus.DebugLevel, args...)
}

func (l *SpotLogger) Info(args ...interface{}) {
	l.Log(logrus.InfoLevel, args...)
}

func (l *SpotLogger) Print(args ...interface{}) {
	l.Log(logrus.InfoLevel, args...)
}

func (l *SpotLogger) Warn(args ...interface{}) {
	l.Log(logrus.WarnLevel, args...)
}

func (l *SpotLogger) Warning(args ...interface{}) {
	l.Warn(args...)
}

func (l *SpotLogger) Error(args ...interface{}) {
	l.Log(logrus.ErrorLevel, args...)
}

func (l *SpotLogger) Fatal(args ...interface{}) {
	l.Log(logrus.FatalLevel, args...)
	l.Exit(1)
}

func (l *SpotLogger) Panic(args ...interface{}) {
	l.Log(logrus.PanicLevel, args...)
}

func (l *SpotLogger) Logln(level logrus.Level, args ...interface{}) {
	l.log(printLogln, level, "", args...)
}

func (l *SpotLogger) Traceln(args ...interface{}) {
	l.Logln(logrus.TraceLevel, args...)
}

func (l *SpotLogger) Debugln(args ...interface{}) {
	l.Logln(logrus.DebugLevel, args...)
}

func (l *SpotLogger) Infoln(args ...interface{}) {
	l.Logln(logrus.InfoLevel, args...)
}

func (l *SpotLogger) Println(args ...interface{}) {
	l.Logln(logrus.InfoLevel, args...)
}

func (l *SpotLogger) Warnln(args ...interface{}) {
	l.Logln(logrus.WarnLevel, args...)
}

func (l *SpotLogger) Warningln(args ...interface{}) {
	l.Warnln(args...)
}

func (l *SpotLogger) Errorln(args ...interface{}) {
	l.Logln(logrus.ErrorLevel, args...)
}

func (l *SpotLogger) Fatalln(args ...interface{}) {
	l.Logln(logrus.FatalLevel, args...)
	l.Exit(1)
}

func (l *SpotLogger) Panicln(args ...interface{}) {
	l.Logln(logrus.PanicLevel, args...)
}
