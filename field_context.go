package zapstackdriver

import (
	"go.uber.org/zap/zapcore"
)

const (
	GcloudErrorEntryKeyContext = "context"
)

type FieldErrorContext struct {
	HttpRequest    *HttpRequest
	User           string
	reportLocation *fieldReportLocation
}

//https://cloud.google.com/error-reporting/reference/rest/v1beta1/ErrorContext#SourceLocation
func (c *FieldErrorContext) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if c.HttpRequest != nil {
		enc.AddObject("httpRequest", c.HttpRequest)
	}

	if c.User != "" {
		enc.AddString("user", c.User)
	}

	if c.reportLocation != nil {
		enc.AddObject("reportLocation", c.reportLocation)
	}

	return nil
}
