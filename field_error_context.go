package zapstackdriver

import (
	"errors"

	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap/zapcore"
)

type FieldErrorContext struct {
	HttpRequest    *HttpRequest
	User           string
	ReportLocation *FieldReportLocation
}

func (c *FieldErrorContext) validate() error {
	var result *multierror.Error

	if c.ReportLocation == nil {
		result = multierror.Append(result, errors.New("report location is required, but was empty"))
	}

	return result.ErrorOrNil()
}

//https://cloud.google.com/error-reporting/reference/rest/v1beta1/ErrorContext#SourceLocation
func (c *FieldErrorContext) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	err := c.validate()

	if err != nil {
		return err
	}

	if c.HttpRequest != nil {
		enc.AddObject("httpRequest", c.HttpRequest)
	}

	if c.User != "" {
		enc.AddString("user", c.User)
	}

	if c.ReportLocation != nil {
		enc.AddObject("reportLocation", c.ReportLocation)
	}

	return nil
}
