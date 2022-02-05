package observer

import (
	"context"
	"fmt"
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	assert.NotNil(t, globalProvider)
	Log("info message")
}

func TestWithLogLevel(t *testing.T) {
	assert.NotNil(t, globalProvider)

	LogWithLevel(LevelDebug, "LevelDebug message")
	LogWithLevel(LevelWarn, "LevelWarn message")
	LogWithLevel(LevelInfo, "LevelInfo message")
	LogWithLevel(LevelNotice, "LevelNotice message")
	LogWithLevel(LevelError, "LevelError message")
	LogWithLevel(LevelAlert, "LevelAlert message")
}

func TestOffAndOn(t *testing.T) {
	assert.NotNil(t, globalProvider)

	DisableLogging()
	_, found := Instance().Find(TypeLogger)
	assert.True(t, found)

	Log("should NOT see this message")

	EnableLogging()
	Log("SHOULD see this message")
}

func TestWithKV(t *testing.T) {
	assert.NotNil(t, globalProvider)

	Log("message with even KVs", "aa", "AA", "bb", "BB", "cc", "CC")
	Log("message with odd KVs", "aa", "AA", "bb")
}

func TestMetering(t *testing.T) {
	assert.NotNil(t, globalProvider)

	Meter(context.Background(), "sample", "aa", "bb")
}

func TestReportError(t *testing.T) {
	assert.NotNil(t, globalProvider)

	// simple exception
	e := fmt.Errorf("an error happened")
	ee := ReportError(e)

	assert.NotNil(t, ee)
	assert.Equal(t, e.Error(), ee.Error())

	// with stacktrace
	ReportError(outer())
}

func inner() error {
	return errors.New("seems we have an error here")
}

func middle() error {
	err := inner()
	if err != nil {
		return err
	}
	return nil
}

func outer() error {
	err := middle()
	if err != nil {
		return err
	}
	return nil
}
