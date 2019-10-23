package zapstackdriver

import (
	"net/http"

	"go.uber.org/zap/zapcore"
)

type FieldHttpRequest struct {
	http.Request
	ResponseStatusCode int
}

func (r FieldHttpRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString(SDKeyErrorContextHttpRequestMethod, r.Method)
	enc.AddString(SDKeyErrorContextHttpRequestUrl, r.URL.String())
	enc.AddString(SDKeyErrorContextHttpRequestRemoteIp, r.RemoteAddr)

	if r.UserAgent() != "" {
		enc.AddString(SDKeyErrorContextHttpRequestUserAgent, r.UserAgent())
	}

	if r.Referer() != "" {
		enc.AddString(SDKeyErrorContextHttpRequestReferrer, r.Referer())
	}

	if http.StatusText(r.ResponseStatusCode) != "" {
		enc.AddInt(SDKeyErrorContextHttpRequestResponseStatusCode, r.ResponseStatusCode)
	}

	return nil
}
