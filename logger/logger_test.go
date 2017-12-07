package logger

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLoggerLevel(t *testing.T) {
	t.Log("Log with respective level (DEBUG, INFO, WARN, ERROR, FATAL, PANIC)")

	var fields logrus.Fields
	var buffer bytes.Buffer

	logrus.SetOutput(&buffer)

	Debug("Lorem Ipsum")

	json.Unmarshal(buffer.Bytes(), &fields)
	if fields["level"] != float64(7) {
		t.Errorf("Wrong level from debug method: %s", fields["level"])
	}

	buffer.Reset()
	Info("Lorem Ipsum")

	json.Unmarshal(buffer.Bytes(), &fields)
	if fields["level"] != float64(6) {
		t.Errorf("Wrong level from info method: %s", fields["level"])
	}

	buffer.Reset()
	Warn("Lorem Ipsum")

	json.Unmarshal(buffer.Bytes(), &fields)
	if fields["level"] != float64(4) {
		t.Errorf("Wrong level from warn method: %s", fields["level"])
	}

	buffer.Reset()
	Error("Lorem Ipsum")

	json.Unmarshal(buffer.Bytes(), &fields)
	if fields["level"] != float64(3) {
		t.Errorf("Wrong level from error method: %s", fields["level"])
	}
}

func TestLoggerCustomFields(t *testing.T) {
	t.Log("Log with custom fields")

	var fields logrus.Fields
	var buffer bytes.Buffer

	logrus.SetOutput(&buffer)
	customFields := WithFields(Fields{"application": "myapp"})

	customFields.Debug("Lorem Ipsum")

	json.Unmarshal(buffer.Bytes(), &fields)
	if fields["application"] != "myapp" {
		t.Errorf("Custom field not logged")
	}

	buffer.Reset()
	fields = logrus.Fields{}
	customFields.Info("Lorem Ipsum")

	json.Unmarshal(buffer.Bytes(), &fields)
	if fields["application"] != "myapp" {
		t.Errorf("Custom field not logged")
	}

	buffer.Reset()
	fields = logrus.Fields{}
	customFields.Warn("Lorem Ipsum")

	json.Unmarshal(buffer.Bytes(), &fields)
	if fields["application"] != "myapp" {
		t.Errorf("Custom field not logged")
	}

	buffer.Reset()
	fields = logrus.Fields{}
	customFields.Error("Lorem Ipsum")

	json.Unmarshal(buffer.Bytes(), &fields)
	if fields["application"] != "myapp" {
		t.Errorf("Custom field not logged")
	}
}
