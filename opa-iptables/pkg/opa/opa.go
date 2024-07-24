package opa

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	Data
}

type Query interface {
	DoQuery(path string, input interface{}) (data []byte, err error)
}

type Data interface {
	PutData(path string, data []byte) error
	GetData(path string) ([]byte, error)
	DeleteData(path string) error
}

type opaClient struct {
	opaEndpoint    string
	authentication string
	client         *http.Client
}

func New(opaEndpoint string, auth string, opaTrustedCAFile string) Client {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	if opaTrustedCAFile != "" {
		if _, err := os.Stat(opaTrustedCAFile); err == nil {
			tlsConfig, err := createTLSConfig(opaTrustedCAFile)
			if err != nil {
				log.Fatalf("Failed to create TLS config: %v", err)
			}
			client.Transport = &http.Transport{
				TLSClientConfig: tlsConfig,
			}
		} else if !os.IsNotExist(err) {
			log.Fatalf("Failed to check CA file: %v", err)
		}
	}
	return &opaClient{opaEndpoint, auth, client}
}

func createTLSConfig(opaTrustedCAFile string) (*tls.Config, error) {
	systemCertPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to load system cert pool: %v", err)
	}
	rootCA, err := ioutil.ReadFile(opaTrustedCAFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read root CA certificate: %v", err)
	}
	if ok := systemCertPool.AppendCertsFromPEM(rootCA); !ok {
		return nil, fmt.Errorf("failed to append root CA certificate")
	}
	return &tls.Config{
		RootCAs: systemCertPool,
	}, nil
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

func (c *opaClient) PutData(path string, data []byte) error {
	url := c.opaEndpoint + fmt.Sprintf(documentEndpointFmt, path)
	_, err := c.do(http.MethodPut, url, data)
	if err != nil {
		return err
	}
	return nil
}

func (c *opaClient) GetData(path string) ([]byte, error) {
	url := c.opaEndpoint + fmt.Sprintf(documentEndpointFmt, path)
	res, err := c.do(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *opaClient) DeleteData(path string) error {
	url := c.opaEndpoint + fmt.Sprintf(documentEndpointFmt, path)
	_, err := c.do(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	return nil
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
