package zapstackdriver

import (
	"go.uber.org/zap/zapcore"
)

type FieldErrorContext struct {
	HTTPRequest      *FieldHTTPRequest
	User             string
	ReportLocation   FieldReportLocation
	SourceReferences []FieldSourceReference
}

//https://cloud.google.com/error-reporting/reference/rest/v1beta1/ErrorContext#SourceLocation
func (c FieldErrorContext) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if c.HTTPRequest != nil {
		enc.AddObject(
			SDKeyErrorContextHTTPRequest,
			c.HTTPRequest,
		)
	}

	if c.User != "" {
		enc.AddString(
			SDKeyErrorContextUser,
			c.User,
		)
	}

	enc.AddObject(
		SDKeyErrorContextReportLocation,
		c.ReportLocation,
	)

	if c.SourceReferences != nil && len(c.SourceReferences) > 0 {
		iterator := func(enc zapcore.ArrayEncoder) error {
			for _, sourceReference := range c.SourceReferences {
				enc.AppendObject(sourceReference)
			}

			return nil
		}

		enc.AddArray(
			SDKeyErrorContextSourceReferences,
			zapcore.ArrayMarshalerFunc(iterator),
		)
	}

	return nil
}
