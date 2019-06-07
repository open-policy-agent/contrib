package opa

import (
	"errors"
)

type Client interface {
	Data
	Policy
}

type Data interface {
	GetData(path string) error
	PutData(path string, data interface{}) error
	PatchData(path string, op string, data interface{}) error
}

type Policy interface {
	InsertPolicy(id string, data []byte) error
	DeletePolicy(id string) error
}

type httpClient struct {
	url string
	authentication string
}

func New(url string, auth string) Client {
	return &httpClient{url,auth}
}

func (c *httpClient) InsertPolicy(id string, data []byte) error {
	return errors.New("Not implemented")
}

func (c *httpClient) DeletePolicy(id string) error {
	return errors.New("Not implemented")
}

func (c *httpClient) GetData(path string) error {
	return errors.New("Not implemented")
}

func (c *httpClient) PatchData(path string,op string, data interface{}) error {
	return errors.New("Not implemented")
}

func (c *httpClient) PutData(path string, data interface{}) error {
	return errors.New("Not implemented")
}