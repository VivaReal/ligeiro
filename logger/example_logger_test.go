package logger_test

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/olxbr/ligeiro/logger"
	"github.com/sirupsen/logrus"
)

func ExampleInfo() {
	var fields logrus.Fields
	var buffer bytes.Buffer
	logrus.SetOutput(&buffer)

	logger.Info("Lorem Ipsum")

	json.Unmarshal(buffer.Bytes(), &fields)

	fmt.Println(fields["short_message"])
	fmt.Println(fields["_environment"])
	fmt.Println(fields["_app_version"])
	// Output:
	// Lorem Ipsum
	// dev
	// detached
}

func ExampleInfof() {
	var fields logrus.Fields
	var buffer bytes.Buffer
	logrus.SetOutput(&buffer)

	logger.WithFields(logger.Fields{"custom": "field"}).
		Infof("Lorem %s", "Ipsum")

	json.Unmarshal(buffer.Bytes(), &fields)

	fmt.Println(fields["short_message"])
	fmt.Println(fields["_environment"])
	fmt.Println(fields["_app_version"])
	fmt.Println(fields["_custom"])
	// Output:
	// Lorem Ipsum
	// dev
	// detached
	// field
}
