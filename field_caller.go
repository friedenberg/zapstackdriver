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
func (f *FieldReportLocation) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("filePath", f.File)
	enc.AddInt("lineNumber", f.Line)
	enc.AddString("functionName", f.FunctionName())
	return nil
}

type FieldSourceLocation struct {
	Caller
}

//https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogEntrySourceLocation
func (f *FieldSourceLocation) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("file", f.File)
	enc.AddString("line", fmt.Sprintf("%d", f.Line))
	enc.AddString("function", f.FunctionName())
	return nil
}
