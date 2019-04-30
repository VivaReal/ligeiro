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
	"bytes"
	"encoding/json"
	"fmt"
	"log/syslog"
	"os"
	"time"

	"github.com/grupozap/ligeiro/envcfg"
	"github.com/sirupsen/logrus"
)

var config = envcfg.LoadBundled()

type entry struct {
	*logrus.Entry
}

const (
	appversionKey  = "app_version"
	environmentKey = "environment"
	fullmessageKey = "full_message"
	levelKey       = "level"
	messageKey     = "short_message"
	specversionKey = "version"
	timestampKey   = "timestamp"
)

type Fields logrus.Fields
type jsonFormatter struct{}

func init() {
	logLevel, err := logrus.ParseLevel(config.Get("logLevel"))
	if err != nil {
		logLevel = logrus.DebugLevel
	}
	logrus.SetLevel(logLevel)

	logrus.SetFormatter(&jsonFormatter{})

	logrus.RegisterExitHandler(func() {
		Info("Application will stop probably due to a OS signal")
	})

	logrus.SetOutput(os.Stdout)
}

func WithFields(fields Fields) *entry {
	fields[environmentKey] = config.Get("environment")
	fields[appversionKey] = config.Get("version")

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
		WithField(levelKey, syslog.LOG_DEBUG).
		WithField(timestampKey, nowMillis()).
		Debug(msg)
}

func (entry *entry) Info(msg interface{}) {
	entry.Entry.
		WithField(levelKey, syslog.LOG_INFO).
		WithField(timestampKey, nowMillis()).
		Info(msg)
}

func (entry *entry) Warn(msg interface{}) {
	entry.Entry.
		WithField(levelKey, syslog.LOG_WARNING).
		WithField(timestampKey, nowMillis()).
		Warn(msg)
}

func (entry *entry) Error(msg interface{}) {
	entry.Entry.
		WithField(levelKey, syslog.LOG_ERR).
		WithField(timestampKey, nowMillis()).
		Error(msg)
}

func (entry *entry) Fatal(msg interface{}) {
	entry.Entry.
		WithField(levelKey, syslog.LOG_CRIT).
		WithField(timestampKey, nowMillis()).
		Fatal(msg)
}

func (entry *entry) Panic(msg interface{}) {
	entry.Entry.
		WithField(levelKey, syslog.LOG_EMERG).
		WithField(timestampKey, nowMillis()).
		Panic(msg)
}

func (e *entry) WithFields(fields Fields) *entry {
	return &entry{e.Entry.WithFields(logrus.Fields(fields))}
}

// Unix timestamp in milliseconds resolution
func nowMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (f *jsonFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(Fields, len(entry.Data)+3)
	for k, v := range entry.Data {
		// GELF spec tells that custom fields includes an `_` as prefix.
		if !(k == fullmessageKey || k == levelKey || k == timestampKey || k == specversionKey) {
			k = "_" + k
		}

		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/sirupsen/logrus/issues/137
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}

	data[messageKey] = entry.Message
	data[specversionKey] = "1.1"

	if _, exists := data[fullmessageKey]; !exists {
		data[fullmessageKey] = ""
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	err := json.NewEncoder(b).Encode(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return b.Bytes(), nil
}
