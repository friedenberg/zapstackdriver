package zapstackdriver_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/friedenberg/zapstackdriver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type jsonObject = map[string]interface{}
type jsonArray = []interface{}
type entry func(l Logger) []Logger
type assertion func(t *testing.T, originalLogger Logger, newLoggers []Logger, actual []jsonObject)

type loggerTestCase struct {
	description      string
	serviceContext   zapstackdriver.FieldServiceContext
	sourceReferences []zapstackdriver.FieldSourceReference
	entry            entry
	assertion        assertion
}

func getLoggerTestCases() []loggerTestCase {
	return []loggerTestCase{
		{
			description: "same logger with Error",
			serviceContext: zapstackdriver.FieldServiceContext{
				Service: "zapstackdriver",
				Version: "1.0",
			},
			sourceReferences: []zapstackdriver.FieldSourceReference{
				{
					Repository: "https://github.com/friedenberg/zapstackdriver",
					RevisionId: "some_sha",
				},
			},
			entry: func(l Logger) []Logger {
				l.Error("simple")
				return []Logger{}
			},
			assertion: func(t *testing.T, originalLogger Logger, newLoggers []Logger, entries []jsonObject) {
				assert.Len(t, entries, 1)
				actual := entries[0]
				assert.Equal(t, "zapstackdriver", actual["serviceContext"].(jsonObject)["service"])
				assert.Equal(t, "1.0", actual["serviceContext"].(jsonObject)["version"])
				assert.Equal(t, zapstackdriver.SDValueErrorType, actual["@type"])
				assert.Equal(t, "simple", strings.Split(actual["message"].(string), "\n")[0])

				context := actual["context"].(jsonObject)
				sourceReferences := context["sourceReferences"].(jsonArray)
				assert.Len(t, sourceReferences, 1)
				assert.Equal(t, "some_sha", sourceReferences[0].(jsonObject)["revisionId"])
				assert.Equal(t, "https://github.com/friedenberg/zapstackdriver", sourceReferences[0].(jsonObject)["repository"])
			},
		},
		{
			description: "different logger with Error",
			serviceContext: zapstackdriver.FieldServiceContext{
				Service: "zapstackdriver",
				Version: "1.0",
			},
			sourceReferences: []zapstackdriver.FieldSourceReference{
				{
					Repository: "https://github.com/friedenberg/zapstackdriver",
					RevisionId: "some_sha",
				},
			},
			entry: func(l Logger) []Logger {
				l = l.With("extra_key", "extra_value")
				l.Error("simple")
				return []Logger{l}
			},
			assertion: func(t *testing.T, originalLogger Logger, newLoggers []Logger, entries []jsonObject) {
				assert.Len(t, entries, 1)
				actual := entries[0]
				assert.Equal(t, "zapstackdriver", actual["serviceContext"].(jsonObject)["service"])
				assert.Equal(t, "1.0", actual["serviceContext"].(jsonObject)["version"])
				assert.Equal(t, zapstackdriver.SDValueErrorType, actual["@type"])
				assert.Equal(t, "simple", strings.Split(actual["message"].(string), "\n")[0])

				context := actual["context"].(jsonObject)
				sourceReferences := context["sourceReferences"].(jsonArray)
				assert.Len(t, sourceReferences, 1)
				assert.Equal(t, "some_sha", sourceReferences[0].(jsonObject)["revisionId"])
				assert.Equal(t, "https://github.com/friedenberg/zapstackdriver", sourceReferences[0].(jsonObject)["repository"])
			},
		},
	}
}

func runLoggerTestCase(t *testing.T, testcase loggerTestCase) {
	temp, err := ioutil.TempFile("", "zapstackdriver-test-logger.*.json")
	require.NoError(t, err, "Failed to create temp file.")
	defer os.Remove(temp.Name())

	loggerConfig := zapstackdriver.NewConfig(testcase.serviceContext)
	loggerConfig.Config.OutputPaths = []string{temp.Name()}
	logger, err := loggerConfig.Build(testcase.sourceReferences)
	require.NoError(t, err, "Failed to create logger")

	loggers := testcase.entry(logger)

	logged, err := ioutil.ReadFile(temp.Name())
	require.NoError(t, err, "Failed to read entry from temp file")

	lines := strings.Split(string(logged), "\n")

	//drop the last line because logs always append a newline and that results in
	//an empty string
	lines = lines[:len(lines)-1]

	logEntries := make([]jsonObject, 0, len(lines))

	for _, line := range lines {
		var result jsonObject
		err = json.Unmarshal([]byte(line), &result)
		require.NoError(t, err, "Failed to unmarshal entry json")
		logEntries = append(logEntries, result)
	}

	testcase.assertion(t, logger, loggers, logEntries)
}

func testCaseRunAdapter(t *testing.T, testcase loggerTestCase, runner func(t *testing.T, testcase loggerTestCase)) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		runner(t, testcase)
	}
}

func TestLogger(t *testing.T) {
	for _, testcase := range getLoggerTestCases() {
		t.Run(testcase.description, testCaseRunAdapter(t, testcase, runLoggerTestCase))
	}
}
