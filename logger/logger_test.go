package logger

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestRequiredFieldByGELPSpec(t *testing.T) {
	t.Log("Comply with GELF 1.1 spec")

	var fields logrus.Fields
	var buffer bytes.Buffer
	logrus.SetOutput(&buffer)

	Info("Lorem Ipsum")
	json.Unmarshal(buffer.Bytes(), &fields)

	if _, exists := fields["full_message"]; !exists {
		t.Error("GELF spec requires `full_message` field to be defined even it's blank")
	}
	if fields["version"] != "1.1" {
		t.Errorf("Current implemented GELF spec version is 1.1, got: %s", fields["version"])
	}
	if _, exists := fields["_version"]; exists {
		t.Error("Field version can't be set as custom")
	}
}

func TestLoggerLevel(t *testing.T) {
	t.Log("Log with respective level (DEBUG, INFO, WARN, ERROR, FATAL, PANIC)")

	var fields logrus.Fields
	var buffer bytes.Buffer

	logrus.SetOutput(&buffer)

	Debug("Lorem Ipsum")

	json.Unmarshal(buffer.Bytes(), &fields)
	if fields["level"] != float64(14) {
		t.Errorf("Wrong level from debug method: %s", fields["level"])
	}

	buffer.Reset()
	Info("Lorem Ipsum")

	json.Unmarshal(buffer.Bytes(), &fields)
	if fields["level"] != float64(13) {
		t.Errorf("Wrong level from info method: %s", fields["level"])
	}

	buffer.Reset()
	Warn("Lorem Ipsum")

	json.Unmarshal(buffer.Bytes(), &fields)
	if fields["level"] != float64(11) {
		t.Errorf("Wrong level from warn method: %s", fields["level"])
	}

	buffer.Reset()
	Error("Lorem Ipsum")

	json.Unmarshal(buffer.Bytes(), &fields)
	if fields["level"] != float64(10) {
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
	if fields["_application"] != "myapp" {
		t.Errorf("Custom field not logged")
	}

	buffer.Reset()
	fields = logrus.Fields{}
	customFields.Info("Lorem Ipsum")

	json.Unmarshal(buffer.Bytes(), &fields)
	if fields["_application"] != "myapp" {
		t.Errorf("Custom field not logged")
	}

	buffer.Reset()
	fields = logrus.Fields{}
	customFields.Warn("Lorem Ipsum")

	json.Unmarshal(buffer.Bytes(), &fields)
	if fields["_application"] != "myapp" {
		t.Errorf("Custom field not logged")
	}

	buffer.Reset()
	fields = logrus.Fields{}
	customFields.Error("Lorem Ipsum")

	json.Unmarshal(buffer.Bytes(), &fields)
	if fields["_application"] != "myapp" {
		t.Errorf("Custom field not logged")
	}
}
