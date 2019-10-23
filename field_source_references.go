package zapstackdriver

import (
	"errors"

	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap/zapcore"
)

type FieldSourceReference struct {
	Repository string
	RevisionId string
}

func (f FieldSourceReference) validate() error {
	var result *multierror.Error

	if f.Repository == "" {
		result = multierror.Append(result, errors.New("repository is required, but was empty"))
	}

	if f.RevisionId == "" {
		result = multierror.Append(result, errors.New("revision id is required, but was empty"))
	}

	return result.ErrorOrNil()
}

//https://cloud.google.com/error-reporting/reference/rest/v1beta1/ErrorContext#SourceLocation
func (f FieldSourceReference) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	err := f.validate()

	if err != nil {
		return err
	}

	enc.AddString(SDKeyErrorContextSourceReferencesRepository, f.Repository)
	enc.AddString(SDKeyErrorContextSourceReferencesRevisionId, f.RevisionId)

	return nil
}
