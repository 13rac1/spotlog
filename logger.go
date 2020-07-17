package spotlog

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// func New() *Logger {
// 	return &Logger{
// 		Out:          os.Stderr,
// 		Formatter:    new(TextFormatter),
// 		Hooks:        make(LevelHooks),
// 		Level:        InfoLevel,
// 		ExitFunc:     os.Exit,
// 		ReportCaller: false,
// 	}
// }

// Logger wraps logrus.Logger to add log storage.
type Logger struct {
	*logrus.Logger
	// minLogLevel is the minimum log level to output.
	minLogLevel logrus.Level

	entries     []storedEntry
	entriesLock sync.Mutex
}

func (l *Logger) alwaysLog(level logrus.Level) bool {
	// Levels have lower values the higher their priority is.
	return level <= l.minLogLevel
}

func (l *Logger) newEntry() *Entry {
	// TODO: Use Pool
	// entry, ok := l.entryPool.Get().(*Entry)
	// if ok {
	// 	return entry
	// }
	return NewEntry(l)
}

func (l *Logger) releaseEntry(entry *Entry) {
	//entry.Data = map[string]interface{}{}
	//l.entryPool.Put(entry)
}

// WithField allocates a new entry and adds a field to it.
// Debug, Print, Info, Warn, Error, Fatal or Panic must be then applied to
// this new returned entry.
// If you want multiple fields, use `WithFields`.
func (l *Logger) WithField(key string, value interface{}) *Entry {
	entry := l.newEntry()
	defer l.releaseEntry(entry)
	return entry.WithField(key, value)
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (l *Logger) WithFields(fields logrus.Fields) *Entry {
	entry := l.newEntry()
	defer l.releaseEntry(entry)
	return entry.WithFields(fields)
}

// Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
func (l *Logger) WithError(err error) *Entry {
	entry := l.newEntry()
	defer l.releaseEntry(entry)
	return entry.WithError(err)
}

// Add a context to the log entry.
func (l *Logger) WithContext(ctx context.Context) *Entry {
	entry := l.newEntry()
	defer l.releaseEntry(entry)
	return entry.WithContext(ctx)
}

// Overrides the time of the log entry.
func (l *Logger) WithTime(t time.Time) *Entry {
	entry := l.newEntry()
	defer l.releaseEntry(entry)
	return entry.WithTime(t)
}

func (l *Logger) log(method printType, level logrus.Level, format string, args ...interface{}) {
	l.entriesLock.Lock()
	defer l.entriesLock.Unlock()

	if l.alwaysLog(level) {
		// Found an important log, print the stored log entries.
		// TODO: Performance: - re-use the logrus.Entry instance.
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

func (l *Logger) Logf(level logrus.Level, format string, args ...interface{}) {
	l.log(printLogf, level, format, args...)
}

func (l *Logger) Tracef(format string, args ...interface{}) {
	l.Logf(logrus.TraceLevel, format, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Logf(logrus.DebugLevel, format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Logf(logrus.InfoLevel, format, args...)
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.Logf(logrus.InfoLevel, format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Logf(logrus.WarnLevel, format, args...)
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.Warnf(format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Logf(logrus.ErrorLevel, format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Logf(logrus.FatalLevel, format, args...)
	l.Exit(1)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.Logf(logrus.PanicLevel, format, args...)
}

func (l *Logger) Log(level logrus.Level, args ...interface{}) {
	l.log(printLog, level, "", args...)
}

func (l *Logger) Trace(args ...interface{}) {
	l.Log(logrus.TraceLevel, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.Log(logrus.DebugLevel, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.Log(logrus.InfoLevel, args...)
}

func (l *Logger) Print(args ...interface{}) {
	l.Log(logrus.InfoLevel, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.Log(logrus.WarnLevel, args...)
}

func (l *Logger) Warning(args ...interface{}) {
	l.Warn(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.Log(logrus.ErrorLevel, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.Log(logrus.FatalLevel, args...)
	l.Exit(1)
}

func (l *Logger) Panic(args ...interface{}) {
	l.Log(logrus.PanicLevel, args...)
}

func (l *Logger) Logln(level logrus.Level, args ...interface{}) {
	l.log(printLogln, level, "", args...)
}

func (l *Logger) Traceln(args ...interface{}) {
	l.Logln(logrus.TraceLevel, args...)
}

func (l *Logger) Debugln(args ...interface{}) {
	l.Logln(logrus.DebugLevel, args...)
}

func (l *Logger) Infoln(args ...interface{}) {
	l.Logln(logrus.InfoLevel, args...)
}

func (l *Logger) Println(args ...interface{}) {
	l.Logln(logrus.InfoLevel, args...)
}

func (l *Logger) Warnln(args ...interface{}) {
	l.Logln(logrus.WarnLevel, args...)
}

func (l *Logger) Warningln(args ...interface{}) {
	l.Warnln(args...)
}

func (l *Logger) Errorln(args ...interface{}) {
	l.Logln(logrus.ErrorLevel, args...)
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.Logln(logrus.FatalLevel, args...)
	l.Exit(1)
}

func (l *Logger) Panicln(args ...interface{}) {
	l.Logln(logrus.PanicLevel, args...)
}
