package clients

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/leonardochaia/vendopunkto/errors"
)

// HTTP wraps http.Client in order to provide error handling
type HTTP interface {
	GetJSON(url string, result interface{}) (*http.Response, error)
	PostJSON(url string, body, result interface{}) (*http.Response, error)
}

// NewHTTPClient creates a new HTTP
func NewHTTPClient() HTTP {
	return &httpImpl{
		client: http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

type httpImpl struct {
	client http.Client
}

func handleResponse(
	op errors.Op,
	path errors.PathName,
	resp *http.Response,
	result interface{},
	err error) (*http.Response, error) {

	if err != nil {
		return resp, errors.E(op, path, err)
	}

	if resp.StatusCode >= 400 {
		return resp, errors.E(op, path, errors.DecodeError(resp))
	}

	if result != nil {
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return resp, errors.E(op, path, err)
		}
	}

	defer resp.Body.Close()

	return resp, nil
}

func (c *httpImpl) GetJSON(url string, result interface{}) (*http.Response, error) {
	const op errors.Op = "http.get"
	path := errors.PathName(url)
	resp, err := c.client.Get(url)
	return handleResponse(op, path, resp, &result, err)
}

func (c *httpImpl) PostJSON(url string, body, result interface{}) (*http.Response, error) {
	const op errors.Op = "http.postJson"
	path := errors.PathName(url)
	params, err := json.Marshal(body)

	if err != nil {
		return nil, errors.E(op, path, errors.Parameters, err)
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(params))
	return handleResponse(op, path, resp, &result, err)
}
