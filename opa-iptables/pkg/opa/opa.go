package opa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	documentEndpointFmt = `/v1/data/%s`
)

// Error contains the standard error fields returned by OPA.
type Error struct {
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Errors  json.RawMessage `json:"errors,omitempty"`
}

func (err *Error) Error() string {
	return fmt.Sprintf("code %v: %v", err.Code, err.Message)
}

type Client interface {
	Query
}

type Query interface {
	DoQuery(path string, input interface{}) (data []byte, err error)
}

type opaClient struct {
	opaEndpoint    string
	authentication string
	client         *http.Client
}

func New(opaEndpoint string, auth string) Client {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	return &opaClient{opaEndpoint, auth, client}
}

func (c *opaClient) DoQuery(path string, input interface{}) (data []byte, err error) {
	url := c.opaEndpoint + fmt.Sprintf(documentEndpointFmt, path)
	d, ok := input.([]byte)
	if !ok {
		return nil, fmt.Errorf("Invalid data; must be []byte")
	}
	res, err := c.do(http.MethodPost, url, d)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *opaClient) do(method, url string, data []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	if c.authentication != "" {
		req.Header.Add("Authorization", "Bearer "+c.authentication)
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = c.handleErrors(res)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(res.Body)
}

func (c *opaClient) handleErrors(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	var err Error
	if err := json.NewDecoder(resp.Body).Decode(&err); err != nil {
		return err
	}
	return &err
}
