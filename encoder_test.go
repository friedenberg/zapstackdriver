package zapstackdriver_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/friedenberg/zapstackdriver"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func TestEncodeEntry(t *testing.T) {
	tests := []struct {
		description string
		entry       zapcore.Entry
		fields      []zap.Field
		expected    string
	}{
		{
			"debug log",
			zapcore.Entry{
				Level: zap.DebugLevel,
			},
			[]zap.Field{},
			"{\"severity\":\"DEBUG\",\"timestamp\":\"0001-01-01T00:00:00.000Z\",\"message\":\"\",\"sourceLocation\":{\"file\":\"\",\"line\":\"0\",\"function\":\"\"}}\n",
		},
		{
			"info log",
			zapcore.Entry{
				Level: zap.InfoLevel,
			},
			[]zap.Field{},
			"{\"severity\":\"INFO\",\"timestamp\":\"0001-01-01T00:00:00.000Z\",\"message\":\"\",\"sourceLocation\":{\"file\":\"\",\"line\":\"0\",\"function\":\"\"}}\n",
		},
		{
			"error log",
			zapcore.Entry{
				Level: zap.ErrorLevel,
			},
			[]zap.Field{},
			"{\"severity\":\"ERROR\",\"timestamp\":\"0001-01-01T00:00:00.000Z\",\"message\":\"\",\"context\":{\"reportLocation\":{\"filePath\":\"\",\"lineNumber\":0,\"functionName\":\"\"}},\"@type\":\"type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent\"}\n",
		},
		{
			"panic log",
			zapcore.Entry{
				Level: zap.PanicLevel,
			},
			[]zap.Field{},
			"{\"severity\":\"ALERT\",\"timestamp\":\"0001-01-01T00:00:00.000Z\",\"message\":\"\",\"context\":{\"reportLocation\":{\"filePath\":\"\",\"lineNumber\":0,\"functionName\":\"\"}},\"@type\":\"type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent\"}\n",
		},
	}

	subject := zapstackdriver.NewEncoder()

	for _, test := range tests {
		t.Run(
			test.description, func(t *testing.T) {
				actual, err := subject.EncodeEntry(test.entry, test.fields)

				assert.Nil(t, err)
				assert.Equal(t, test.expected, actual.String())
			},
		)
	}
}

func TestEncoderPersistsInClonedLogger(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping")
	}

	var b bytes.Buffer
	logger := zaptest.NewLogger(t)
	logger = logger.WithOptions(
		zapstackdriver.NewCore(&b),
		zapstackdriver.WithContext(
			"0",
			"1",
		),
	)

	logger.Sugar().Errorw("foo", "bar", "baz", "qux", "quux")

	assert.Contains(t, b.String(), "@type")
	assert.Contains(t, b.String(), "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent")
	assert.Contains(t, b.String(), "\"message\":\"foo\"")
	assert.Contains(t, b.String(), "\"bar\":\"baz\"")
	assert.Contains(t, b.String(), "\"qux\":\"quux\"")
	assert.Contains(t, b.String(), "\"serviceContext\":{\"service\":\"0\",\"version\":\"1\"}")
	assert.Contains(t, b.String(), fmt.Sprintf("\"severity\":\"%s\"", "ERROR"))

	anotherLogger := logger.Sugar().With("anotherkey", "anothervalue")
	anotherLogger.Errorw("foo", "bar", "baz", "qux", "quux")

	assert.Contains(t, b.String(), "\"anotherkey\":\"anothervalue\"")
	assert.Contains(t, b.String(), "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent")
	assert.Contains(t, b.String(), "@type")
	assert.Contains(t, b.String(), "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent")
	assert.Contains(t, b.String(), "\"message\":\"foo\"")
	assert.Contains(t, b.String(), "\"bar\":\"baz\"")
	assert.Contains(t, b.String(), "\"qux\":\"quux\"")
	assert.Contains(t, b.String(), "\"serviceContext\":{\"service\":\"0\",\"version\":\"1\"}")
	assert.Contains(t, b.String(), fmt.Sprintf("\"severity\":\"%s\"", "ERROR"))
}
