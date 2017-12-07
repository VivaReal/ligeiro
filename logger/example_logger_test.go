package logger_test

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/vivareal/ligeiro/logger"
)

func ExampleInfo() {
	var fields logrus.Fields
	var buffer bytes.Buffer
	logrus.SetOutput(&buffer)

	logger.Info("Lorem Ipsum")

	json.Unmarshal(buffer.Bytes(), &fields)

	fmt.Println(fields["full_message"])
	fmt.Println(fields["environment"])
	fmt.Println(fields["version"])
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

	fmt.Println(fields["full_message"])
	fmt.Println(fields["environment"])
	fmt.Println(fields["version"])
	fmt.Println(fields["custom"])
	// Output:
	// Lorem Ipsum
	// dev
	// detached
	// field
}
