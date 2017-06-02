package logging_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/RichardKnop/logging"
	"github.com/stretchr/testify/assert"
)

func TestDefaultFormatter(t *testing.T) {
	var (
		out, errOut = bytes.NewBuffer([]byte{}), bytes.NewBuffer([]byte{})
		logger      = logging.New(out, errOut, nil)
		now         time.Time
		actual      []byte
		expected    string
		err         error
	)

	// Test logger.Info
	now = time.Now()
	logger[logging.INFO].Print("Test logger.Print")
	actual, err = ioutil.ReadAll(out)
	if err != nil {
		log.Fatal(err)
	}
	expected = fmt.Sprintf(
		"INFO: %s default_formatter_test.go:27 Test logger.Print\n",
		now.Format("2006/01/02 15:04:05"),
	)
	assert.Equal(t, expected, string(actual))

	// Test logger.Infof
	now = time.Now()
	logger[logging.INFO].Printf("Test %s.%s", "logger", "Printf")
	actual, err = ioutil.ReadAll(out)
	if err != nil {
		log.Fatal(err)
	}
	expected = fmt.Sprintf(
		"INFO: %s default_formatter_test.go:40 Test logger.Printf\n",
		now.Format("2006/01/02 15:04:05"),
	)
	assert.Equal(t, expected, string(actual))

	// Test logger.Warning
	now = time.Now()
	logger[logging.WARNING].Print("Test logger.Print")
	actual, err = ioutil.ReadAll(out)
	if err != nil {
		log.Fatal(err)
	}
	expected = fmt.Sprintf(
		"WARNING: %s default_formatter_test.go:53 Test logger.Print\n",
		now.Format("2006/01/02 15:04:05"),
	)
	assert.Equal(t, expected, string(actual))

	// Test logger.Warningf
	now = time.Now()
	logger[logging.WARNING].Printf("Test %s.%s", "logger", "Printf")
	actual, err = ioutil.ReadAll(out)
	if err != nil {
		log.Fatal(err)
	}
	expected = fmt.Sprintf(
		"WARNING: %s default_formatter_test.go:66 Test logger.Printf\n",
		now.Format("2006/01/02 15:04:05"),
	)
	assert.Equal(t, expected, string(actual))

	// Test logger.Error
	now = time.Now()
	logger[logging.ERROR].Print("Test logger.Print")
	actual, err = ioutil.ReadAll(errOut)
	if err != nil {
		log.Fatal(err)
	}
	expected = fmt.Sprintf(
		"ERROR: %s default_formatter_test.go:79 Test logger.Print\n",
		now.Format("2006/01/02 15:04:05"),
	)
	assert.Equal(t, expected, string(actual))

	// Test logger.Errorf
	now = time.Now()
	logger[logging.ERROR].Printf("Test %s.%s", "logger", "Printf")
	actual, err = ioutil.ReadAll(errOut)
	if err != nil {
		log.Fatal(err)
	}
	expected = fmt.Sprintf(
		"ERROR: %s default_formatter_test.go:92 Test logger.Printf\n",
		now.Format("2006/01/02 15:04:05"),
	)
	assert.Equal(t, expected, string(actual))

	// Test logger.Fatal
	now = time.Now()
	logger[logging.FATAL].Print("Test logger.Print")
	actual, err = ioutil.ReadAll(errOut)
	if err != nil {
		log.Fatal(err)
	}
	expected = fmt.Sprintf(
		"FATAL: %s default_formatter_test.go:105 Test logger.Print\n",
		now.Format("2006/01/02 15:04:05"),
	)
	assert.Equal(t, expected, string(actual))

	// Test logger.Fatalf
	now = time.Now()
	logger[logging.FATAL].Printf("Test %s.%s", "logger", "Printf")
	actual, err = ioutil.ReadAll(errOut)
	if err != nil {
		log.Fatal(err)
	}
	expected = fmt.Sprintf(
		"FATAL: %s default_formatter_test.go:118 Test logger.Printf\n",
		now.Format("2006/01/02 15:04:05"),
	)
	assert.Equal(t, expected, string(actual))
}
