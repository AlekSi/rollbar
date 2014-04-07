package rollbar

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"
)

var platform string

func init() {
	platform = runtime.GOOS + " " + runtime.GOARCH
}

type Person struct {
	ID       string `json:"id"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

type Payload struct {
	Error       error
	Message     string
	MessageMeta map[string]interface{}
	Level       string // one of: critical, error, warning, info, debug
	Timestamp   time.Time
	Context     string
	Request     *http.Request
	Person      *Person
}

func (p *Payload) toMap() map[string]interface{} {
	body := make(map[string]interface{})

	switch p.Error {
	case nil:
		message := map[string]interface{}{
			"body": p.Message,
		}
		for k, v := range p.MessageMeta {
			message[k] = v
		}
		body["message"] = message

	default:
		trace := map[string]interface{}{
			"frames": stack(),
			"exception": map[string]interface{}{
				"class":   reflect.TypeOf(p.Error).String(),
				"message": p.Error.Error(),
			},
		}
		body["trace"] = trace
	}

	data := map[string]interface{}{
		"body":     body,
		"platform": platform,
		"language": "go",
	}

	if p.Level != "" {
		data["level"] = p.Level
	}

	if !p.Timestamp.IsZero() {
		data["timestamp"] = p.Timestamp.Unix()
	}

	if p.Context != "" {
		data["context"] = p.Context
	}

	if p.Request != nil {
		url := p.Request.RequestURI
		if url == "" {
			url = p.Request.URL.String()
		}

		headers := make(map[string]interface{}, len(p.Request.Header))
		for k, v := range p.Request.Header {
			headers[k] = strings.Join(v, ", ")
		}

		data["request"] = map[string]interface{}{
			"url":     url,
			"method":  p.Request.Method,
			"headers": headers,
		}
	}

	if p.Person != nil {
		data["person"] = p.Person
	}

	return map[string]interface{}{
		"data": data,
	}
}
