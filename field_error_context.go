package zapstackdriver

import (
	"go.uber.org/zap/zapcore"
)

type FieldErrorContext struct {
	HttpRequest      *HttpRequest
	User             string
	ReportLocation   FieldReportLocation
	SourceReferences []FieldSourceReference
}

//https://cloud.google.com/error-reporting/reference/rest/v1beta1/ErrorContext#SourceLocation
func (c *FieldErrorContext) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if c.HttpRequest != nil {
		enc.AddObject("httpRequest", c.HttpRequest)
	}

	if c.User != "" {
		enc.AddString("user", c.User)
	}

	enc.AddObject("reportLocation", c.ReportLocation)

	if c.SourceReferences != nil && len(c.SourceReferences) > 0 {
		iterator := func(enc zapcore.ArrayEncoder) error {
			for _, sourceReference := range c.SourceReferences {
				enc.AppendObject(sourceReference)
			}

			return nil
		}

		enc.AddArray("sourceReferences", zapcore.ArrayMarshalerFunc(iterator))
	}

	return nil
}