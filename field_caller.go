package zapstackdriver

import (
	"fmt"
	"runtime"

	"go.uber.org/zap/zapcore"
)

type Caller struct {
	zapcore.EntryCaller
}

func (c Caller) FunctionName() string {
	return runtime.FuncForPC(c.PC).Name()
}

type FieldReportLocation struct {
	Caller
}

//https://cloud.google.com/error-reporting/reference/rest/v1beta1/ErrorContext#SourceLocation
func (f FieldReportLocation) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString(SDKeyErrorContextReportLocationFilePath, f.File)
	enc.AddInt(SDKeyErrorContextReportLineNumber, f.Line)
	enc.AddString(SDKeyErrorContextReportFunctionName, f.FunctionName())
	return nil
}

type FieldSourceLocation struct {
	Caller
}

//https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogEntrySourceLocation
func (f FieldSourceLocation) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString(SDKeySourceLocationFile, f.File)
	enc.AddString(SDKeySourceLocationLine, fmt.Sprintf("%d", f.Line))
	enc.AddString(SDKeySourceLocationFunction, f.FunctionName())
	return nil
}
