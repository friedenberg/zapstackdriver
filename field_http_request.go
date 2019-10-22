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
	enc.AddString("method", r.Method)
	enc.AddString("url", r.URL.String())
	enc.AddString("remoteIp", r.RemoteAddr)

	if r.UserAgent() != "" {
		enc.AddString("userAgent", r.UserAgent())
	}

	if r.Referer() != "" {
		enc.AddString("referrer", r.Referer())
	}

	if http.StatusText(r.ResponseStatusCode) != "" {
		enc.AddInt("responseStatusCode", r.ResponseStatusCode)
	}

	return nil
}
