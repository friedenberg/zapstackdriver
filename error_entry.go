package zapstackdriver

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	gcloudErrorEntryKeyType        = "@type"
	gcloudErrorEntryValueErrorType = "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent"
)

type errorEntry struct {
	entry  zapcore.Entry
	fields []zapcore.Field
}

func MakeErrorEntry(entry zapcore.Entry, fields []zapcore.Field) (zapcore.Entry, []zapcore.Field) {
	hasContext := false

	reportLocation := &fieldReportLocation{caller{entry.Caller}}

	for _, field := range fields {
		if field.Key == GcloudErrorEntryKeyContext {
			switch context := field.Interface.(type) {
			case *FieldErrorContext:
				hasContext = true
				context.reportLocation = reportLocation
			default:
				fields = append(fields, zap.String("contextError", "error entry field has incorrect type"))
			}
		}
	}

	if !hasContext {
		fieldValue := &FieldErrorContext{
			reportLocation: reportLocation,
		}

		fields = append(fields, zap.Object(GcloudErrorEntryKeyContext, fieldValue))
	}

	typeField := zap.String(gcloudErrorEntryKeyType, gcloudErrorEntryValueErrorType)

	fields = append(fields, typeField)

	return entry, fields
}
