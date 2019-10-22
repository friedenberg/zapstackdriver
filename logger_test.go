package zapstackdriver_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/friedenberg/zapstackdriver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type jsonObject = map[string]interface{}
type entry func(l Logger) Logger
type assertion func(t *testing.T, actual jsonObject)

func TestLogger(t *testing.T) {
	for _, testcase := range []struct {
		description      string
		serviceContext   zapstackdriver.FieldServiceContext
		sourceReferences []zapstackdriver.FieldSourceReference
		entry            entry
		assertion        assertion
	}{
		{
			description: "simple logger",
			serviceContext: zapstackdriver.FieldServiceContext{
				Service: "zapstackdriver",
				Version: "1.0",
			},
			sourceReferences: []zapstackdriver.FieldSourceReference{
				zapstackdriver.FieldSourceReference{
					Repository: "https://github.com/friedenberg/zapstackdriver",
					RevisionId: "some_sha",
				},
			},
			entry: func(l Logger) Logger {
				l.Errorw("simple")
				return l
			},
			assertion: func(t *testing.T, actual jsonObject) {
				assert.Equal(t, "zapstackdriver", actual["serviceContext"].(jsonObject)["service"])
				assert.Equal(t, "1.0", actual["serviceContext"].(jsonObject)["version"])
				assert.Equal(t, "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent", actual["@type"])
			},
		},
	} {
		t.Run(
			testcase.description,
			func(t *testing.T) {
				temp, err := ioutil.TempFile("", "zapstackdriver-test-logger")
				require.NoError(t, err, "Failed to create temp file.")
				defer os.Remove(temp.Name())

				loggerConfig := zapstackdriver.NewConfig(testcase.serviceContext)
				loggerConfig.Config.OutputPaths = []string{temp.Name()}
				logger, err := loggerConfig.Build(testcase.sourceReferences)
				require.NoError(t, err, "Failed to create logger")

				testcase.entry(logger)

				logged, err := ioutil.ReadFile(temp.Name())
				require.NoError(t, err, "Failed to read entry from temp file")

				var result jsonObject
				err = json.Unmarshal(logged, &result)
				require.NoError(t, err, "Failed to unmarshal entry json")

				testcase.assertion(t, result)
			},
		)
	}
}
