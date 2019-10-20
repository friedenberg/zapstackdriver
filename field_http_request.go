package zapstackdriver

import (
	"errors"
	"net/http"

	"go.uber.org/zap/zapcore"
)

type HttpRequest struct {
	*http.Request
	ResponseStatusCode int
}

func (r HttpRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if r.Request == nil {
		return errors.New("http request context is nil")
	}

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
