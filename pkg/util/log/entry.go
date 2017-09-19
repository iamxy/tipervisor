package log

import (
	"github.com/sirupsen/logrus"
)

// Entry type
type Entry logrus.Entry

// Returns the string representation from the reader and ultimately the
// formatter.
func (entry *Entry) String() (string, error) {
	l := (*logrus.Entry)(entry)
	return l.String()
}

// WithError adds  an error as single field (using the key defined in ErrorKey) to the Entry.
func (entry *Entry) WithError(err error) *Entry {
	l := (*logrus.Entry)(entry)
	return (*Entry)(l.WithError(err))
}

// WithField adds a single field to the Entry.
func (entry *Entry) WithField(key string, value interface{}) *Entry {
	l := (*logrus.Entry)(entry)
	return (*Entry)(l.WithField(key, value))
}

// WithFields adds a map of fields to the Entry.
func (entry *Entry) WithFields(fields Fields) *Entry {
	f := logrus.Fields(fields)
	l := (*logrus.Entry)(entry)
	return (*Entry)(l.WithFields(f))
}

// Debug logs a message at level Debug on the logrus logger.
func (entry *Entry) Debug(args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Debug(args...)
}

// Print logs a message at level Info on the logrus logger.
func (entry *Entry) Print(args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Print(args...)
}

// Info logs a message at level Info on the logrus logger.
func (entry *Entry) Info(args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Info(args...)
}

// Warn logs a message at level Warn on the logrus logger.
func (entry *Entry) Warn(args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Warn(args...)
}

// Warning logs a message at level Warn on the logrus logger.
func (entry *Entry) Warning(args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Warning(args...)
}

// Error logs a message at level Error on the logrus logger.
func (entry *Entry) Error(args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Error(args...)
}

// Fatal logs a message at level Fatal on the logrus logger.
func (entry *Entry) Fatal(args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Fatal(args...)
}

// Panic logs a message at level Panic on the logrus logger.
func (entry *Entry) Panic(args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Panic(args...)
}

// Entry Printf family functions

// Debugf logs a message at level Debug on the logrus logger.
func (entry *Entry) Debugf(format string, args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Debugf(format, args...)
}

// Infof logs a message at level Info on the logrus logger.
func (entry *Entry) Infof(format string, args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Infof(format, args...)
}

// Printf logs a message at level Info on the logrus logger.
func (entry *Entry) Printf(format string, args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Printf(format, args...)
}

// Warnf logs a message at level Warn on the logrus logger.
func (entry *Entry) Warnf(format string, args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Warnf(format, args...)
}

// Warningf logs a message at level Warn on the logrus logger.
func (entry *Entry) Warningf(format string, args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Warningf(format, args...)
}

// Errorf logs a message at level Error on the logrus logger.
func (entry *Entry) Errorf(format string, args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Errorf(format, args...)
}

// Fatalf logs a message at level Fatal on the logrus logger.
func (entry *Entry) Fatalf(format string, args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Fatalf(format, args...)
}

// Panicf logs a message at level Panic on the logrus logger.
func (entry *Entry) Panicf(format string, args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Panicf(format, args...)
}

// Entry Println family functions

// Debugln logs a message at level Debug on the logrus logger.
func (entry *Entry) Debugln(args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Debugln(args...)
}

// Infoln logs a message at level Info on the logrus logger.
func (entry *Entry) Infoln(args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Infoln(args...)
}

// Println logs a message at level Info on the logrus logger.
func (entry *Entry) Println(args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Println(args...)
}

// Warnln logs a message at level Warn on the logrus logger.
func (entry *Entry) Warnln(args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Warnln(args...)
}

// Warningln logs a message at level Warn on the logrus logger.
func (entry *Entry) Warningln(args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Warningln(args...)
}

// Errorln logs a message at level Error on the logrus logger.
func (entry *Entry) Errorln(args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Errorln(args...)
}

// Fatalln logs a message at level Fatal on the logrus logger.
func (entry *Entry) Fatalln(args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Fatalln(args...)
}

// Panicln logs a message at level Panic on the logrus logger.
func (entry *Entry) Panicln(args ...interface{}) {
	l := (*logrus.Entry)(entry)
	l.Panicln(args...)
}
