package rollbar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func fromJson(s string) map[string]interface{} {
	var m map[string]interface{}
	d := json.NewDecoder(bytes.NewBufferString(s))
	err := d.Decode(&m)
	if err != nil {
		panic(err)
	}
	return m
}

func TestPayload(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://rollbar.com", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Cache-Control", "no-store")

	p := &Payload{
		Message: "Hello!",
		Level:   "error",
		Context: "rollbar#test",
		Request: req,
	}

	actual := p.toMap()
	expected := fromJson(fmt.Sprintf(`{
		"data": {
			"body": {
				"message": {
					"body": "Hello!"
				}
			},
			"level": "error",
			"platform": %q,
			"language": "go",
			"context": "rollbar#test",
			"request": {
				"method": "GET",
				"url": "http://rollbar.com",
				"headers": {
					"Content-Type": "application/json",
					"Cache-Control": "no-cache, no-store"
				}
			}
		}
	}`, platform))
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\n%#v\n!=\n%#v", actual, expected)
	}
}
