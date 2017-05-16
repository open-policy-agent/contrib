// Package clair implements helpers to interact with the Clair API.
package clair

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Clair struct {
	baseURL string
}

func New(url string) *Clair {
	return &Clair{
		baseURL: url,
	}
}

type IndexProps struct {
	Name       string
	ParentName string `json:",omitempty"`
	Path       string
	Format     string
	Headers    map[string]string
	Timeout    time.Duration `json:"-"`
}

func (cl *Clair) Index(props IndexProps) error {

	var buf bytes.Buffer

	body := struct {
		Layer IndexProps
	}{props}

	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", cl.baseURL+"/layers", &buf)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	client := http.Client{
		Timeout: props.Timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return cl.decodeError(resp)
	}

	return nil
}

func (cl *Clair) Layer(layer string) (map[string]interface{}, error) {

	req, err := http.NewRequest("GET", cl.baseURL+"/layers/"+layer+"?vulnerabilities", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, cl.decodeError(resp)
	}

	var result map[string]interface{}
	return result, json.NewDecoder(resp.Body).Decode(&result)
}

func (cl *Clair) decodeError(resp *http.Response) error {
	if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		return fmt.Errorf("unexpected status code: %v", resp.Status)
	}
	var body struct {
		Error struct {
			Message string
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return err
	}
	return fmt.Errorf("%v (%v)", body.Error.Message, resp.Status)
}
