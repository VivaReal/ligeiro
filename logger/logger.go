// Wrapper of logrus logging infrastructure.
//
// This package contains all the utilities to make logrus logging library attend
// all requirements of VivaReal logging conventions described here https://github.com/VivaReal/platform/blob/master/Documentation/Logs/format.md
//
// It adds common fields as application, environment, process, product and version and provides a
// simple API:
//
package logger

import (
	"fmt"
	"log/syslog"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/grupozap/ligeiro/envcfg"
)

var config = envcfg.LoadBundled()

type entry struct {
	*logrus.Entry
}

type Fields logrus.Fields

func init() {
	logrus.SetLevel(levelFromString(config.Get("logLevel")))

	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg:   "full_message",
			logrus.FieldKeyLevel: "level_name",
		},
	})

	logrus.RegisterExitHandler(func() {
		Info("Application will stop probably due to a OS signal")
	})

	logrus.SetOutput(os.Stdout)
}

func WithFields(fields Fields) *entry {
	fields["environment"] = config.Get("environment")
	fields["version"] = config.Get("version")

	return &entry{logrus.WithFields(logrus.Fields(fields))}
}

func Debug(msg interface{}) {
	WithFields(Fields{}).Debug(msg)
}

func Debugf(format string, args ...interface{}) {
	WithFields(Fields{}).Debug(fmt.Sprintf(format, args...))
}

func Info(msg interface{}) {
	WithFields(Fields{}).Info(msg)
}

func Infof(format string, args ...interface{}) {
	WithFields(Fields{}).Info(fmt.Sprintf(format, args...))
}

func Warn(msg interface{}) {
	WithFields(Fields{}).Warn(msg)
}

func Warnf(format string, args ...interface{}) {
	WithFields(Fields{}).Warn(fmt.Sprintf(format, args...))
}

func Error(msg interface{}) {
	WithFields(Fields{}).Error(msg)
}

func Errorf(format string, args ...interface{}) {
	WithFields(Fields{}).Error(fmt.Sprintf(format, args...))
}

func Fatal(msg interface{}) {
	WithFields(Fields{}).Fatal(msg)
}

func Fatalf(format string, args ...interface{}) {
	WithFields(Fields{}).Fatal(fmt.Sprintf(format, args...))
}

func Panic(msg interface{}) {
	WithFields(Fields{}).Panic(msg)
}

func Panicf(format string, args ...interface{}) {
	WithFields(Fields{}).Panic(fmt.Sprintf(format, args...))
}

func (entry *entry) Debug(msg interface{}) {
	entry.Entry.
		WithField("level", syslog.LOG_DEBUG).
		WithField("timestamp", nowMillis()).
		Debug(msg)
}

func (entry *entry) Info(msg interface{}) {
	entry.Entry.
		WithField("level", syslog.LOG_INFO).
		WithField("timestamp", nowMillis()).
		Info(msg)
}

func (entry *entry) Warn(msg interface{}) {
	entry.Entry.
		WithField("level", syslog.LOG_WARNING).
		WithField("timestamp", nowMillis()).
		Warn(msg)
}

func (entry *entry) Error(msg interface{}) {
	entry.Entry.
		WithField("level", syslog.LOG_ERR).
		WithField("timestamp", nowMillis()).
		Error(msg)
}

func (entry *entry) Fatal(msg interface{}) {
	entry.Entry.
		WithField("level", syslog.LOG_CRIT).
		WithField("timestamp", nowMillis()).
		Fatal(msg)
}

func (entry *entry) Panic(msg interface{}) {
	entry.Entry.
		WithField("level", syslog.LOG_EMERG).
		WithField("timestamp", nowMillis()).
		Panic(msg)
}

func (e *entry) WithFields(fields Fields) *entry {
	return &entry{e.Entry.WithFields(logrus.Fields(fields))}
}

// Convert the level string to a logrusrus Level. E.g. "panic" becomes "PanicLevel".
func levelFromString(level string) logrus.Level {
	switch level {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.DebugLevel
	}
}

// Unix timestamp in milliseconds resolution
func nowMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
