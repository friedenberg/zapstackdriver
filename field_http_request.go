package zapstackdriver

import (
	"net/http"

	"go.uber.org/zap/zapcore"
)

type FieldHTTPRequest struct {
	http.Request
	ResponseStatusCode int
}

func (r FieldHTTPRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString(SDKeyErrorContextHTTPRequestMethod, r.Method)
	enc.AddString(SDKeyErrorContextHTTPRequestURL, r.URL.String())
	enc.AddString(SDKeyErrorContextHTTPRequestRemoteIP, r.RemoteAddr)

	if r.UserAgent() != "" {
		enc.AddString(SDKeyErrorContextHTTPRequestUserAgent, r.UserAgent())
	}

	if r.Referer() != "" {
		enc.AddString(SDKeyErrorContextHTTPRequestReferrer, r.Referer())
	}

	if http.StatusText(r.ResponseStatusCode) != "" {
		enc.AddInt(SDKeyErrorContextHTTPRequestResponseStatusCode, r.ResponseStatusCode)
	}

	return nil
}
