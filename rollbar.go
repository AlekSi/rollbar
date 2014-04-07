package rollbar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

const (
	httpURL  = "http://api.rollbar.com/api/1/item/"
	httpsURL = "https://api.rollbar.com/api/1/item/"
)

type Client struct {
	Token       string
	Environment string
	HTTPClient  *http.Client
	UseHTTP     bool
}

func (c *Client) http() *http.Client {
	if c.HTTPClient == nil {
		return http.DefaultClient
	}
	return c.HTTPClient
}

func (c *Client) url() string {
	if c.UseHTTP {
		return httpURL
	}
	return httpsURL
}

func (c *Client) Post(payload *Payload) (err error) {
	defer func() {
		if err != nil {
			if nerr, ok := err.(net.Error); ok {
				err = nerr
			} else {
				err = FatalError{err}
			}
		}
	}()

	m := payload.toMap()
	m["access_token"] = c.Token
	if c.Environment != "" {
		m["data"].(map[string]interface{})["environment"] = c.Environment
	}

	b, err := json.Marshal(m)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", c.url(), bytes.NewReader(b))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.http().Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()

	if (res.StatusCode / 100) != 2 {
		b, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return
		}

		err = fmt.Errorf("%s", b)
	}

	return
}

func (c *Client) PostMessage(message string) error {
	return c.Post(&Payload{Message: message})
}

func (c *Client) PostRequestMessage(req *http.Request, message string) error {
	return c.Post(&Payload{Message: message, Request: req})
}

func (c *Client) PostError(err error) error {
	return c.Post(&Payload{Error: err})
}

func (c *Client) PostRequestError(req *http.Request, err error) error {
	return c.Post(&Payload{Error: err, Request: req})
}
