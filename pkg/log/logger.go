package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

// Fields type, used to pass to `WithFields`.
type Fields logrus.Fields

var (
	// L keeps the global logrus logger
	L = logrus.NewEntry(logrus.StandardLogger())
)

// SetOutput sets the logrus logger output.
func SetOutput(out io.Writer) {
	logrus.SetOutput(out)
}

// SetFormatter sets the logrus logger formatter.
func SetFormatter(formatter logrus.Formatter) {
	logrus.SetFormatter(formatter)
}

// SetLevel sets the logrus logger level.
func SetLevel(level logrus.Level) {
	logrus.SetLevel(level)
}

// GetLevel returns the logrus logger level.
func GetLevel() logrus.Level {
	return logrus.GetLevel()
}

// AddHook adds a hook to the logrus logger hooks.
func AddHook(hook logrus.Hook) {
	logrus.AddHook(hook)
}

// SetRootFields sets the root fields to logrus logger. Note that
// calling this is NOT goroutine safe
func SetRootFields(fields Fields) {
	L = logrus.WithFields(logrus.Fields(fields))
}

// WithError creates an entry from the logrus logger and adds an error to it,
// using the value defined in ErrorKey as key.
func WithError(err error) *Entry {
	return (*Entry)(L.WithField(logrus.ErrorKey, err))
}

// WithField creates an entry from the logrus logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *Entry {
	return (*Entry)(L.WithField(key, value))
}

// WithFields creates an entry from the logrus logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields Fields) *Entry {
	return (*Entry)(L.WithFields(logrus.Fields(fields)))
}

// Debug logs a message at level Debug on the logrus logger.
func Debug(args ...interface{}) {
	L.Debug(args...)
}

// Print logs a message at level Info on the logrus logger.
func Print(args ...interface{}) {
	L.Print(args...)
}

// Info logs a message at level Info on the logrus logger.
func Info(args ...interface{}) {
	L.Info(args...)
}

// Warn logs a message at level Warn on the logrus logger.
func Warn(args ...interface{}) {
	L.Warn(args...)
}

// Warning logs a message at level Warn on the logrus logger.
func Warning(args ...interface{}) {
	L.Warning(args...)
}

// Error logs a message at level Error on the logrus logger.
func Error(args ...interface{}) {
	L.Error(args...)
}

// Panic logs a message at level Panic on the logrus logger.
func Panic(args ...interface{}) {
	L.Panic(args...)
}

// Fatal logs a message at level Fatal on the logrus logger.
func Fatal(args ...interface{}) {
	L.Fatal(args...)
}

// Debugf logs a message at level Debug on the logrus logger.
func Debugf(format string, args ...interface{}) {
	L.Debugf(format, args...)
}

// Printf logs a message at level Info on the logrus logger.
func Printf(format string, args ...interface{}) {
	L.Printf(format, args...)
}

// Infof logs a message at level Info on the logrus logger.
func Infof(format string, args ...interface{}) {
	L.Infof(format, args...)
}

// Warnf logs a message at level Warn on the logrus logger.
func Warnf(format string, args ...interface{}) {
	L.Warnf(format, args...)
}

// Warningf logs a message at level Warn on the logrus logger.
func Warningf(format string, args ...interface{}) {
	L.Warningf(format, args...)
}

// Errorf logs a message at level Error on the logrus logger.
func Errorf(format string, args ...interface{}) {
	L.Errorf(format, args...)
}

// Panicf logs a message at level Panic on the logrus logger.
func Panicf(format string, args ...interface{}) {
	L.Panicf(format, args...)
}

// Fatalf logs a message at level Fatal on the logrus logger.
func Fatalf(format string, args ...interface{}) {
	L.Fatalf(format, args...)
}

// Debugln logs a message at level Debug on the logrus logger.
func Debugln(args ...interface{}) {
	L.Debugln(args...)
}

// Println logs a message at level Info on the logrus logger.
func Println(args ...interface{}) {
	L.Println(args...)
}

// Infoln logs a message at level Info on the logrus logger.
func Infoln(args ...interface{}) {
	L.Infoln(args...)
}

// Warnln logs a message at level Warn on the logrus logger.
func Warnln(args ...interface{}) {
	L.Warnln(args...)
}

// Warningln logs a message at level Warn on the logrus logger.
func Warningln(args ...interface{}) {
	L.Warningln(args...)
}

// Errorln logs a message at level Error on the logrus logger.
func Errorln(args ...interface{}) {
	L.Errorln(args...)
}

// Panicln logs a message at level Panic on the logrus logger.
func Panicln(args ...interface{}) {
	L.Panicln(args...)
}

// Fatalln logs a message at level Fatal on the logrus logger.
func Fatalln(args ...interface{}) {
	L.Fatalln(args...)
}
