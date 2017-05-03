// +build darwin linux

// Copyright 2017 The OPA Authors. All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log/syslog"
	"net/http"
	"runtime"
	"strings"
)

const (
	defaultUrl      = "http://localhost:8181"
	defaultHostFile = "/etc/host_identity.json"
)

// AuthResult is the result of the authentcate function.
type AuthResult int

const (
	// AuthError is a failure.
	AuthError AuthResult = iota
	// AuthSuccess is a success.
	AuthSuccess
)

// Input to the OPA policy
type authzPolicyInput struct {
	Input struct {
		User         string      `json:"user"`
		HostIdentity interface{} `json:"host_identity"`
	} `json:"input"`
}

// The response from OPA is expected to have these fields.
// In other words, the policy should bind data to 'allow' and 'errors'
type authzPolicyResult struct {
	Allow  bool     `json:"allow"`
	Errors []string `json:"errors"`
}

// Response format from the OPA policy evaluation API
type authzPolicyResponse struct {
	Result authzPolicyResult `json:"result,omitempty"`
}

func pamLog(format string, args ...interface{}) {
	l, err := syslog.New(syslog.LOG_AUTH|syslog.LOG_WARNING, "pam-authz")
	if err != nil {
		return
	}
	l.Warning(fmt.Sprintf(format, args...))
}

// Call over HTTP to OPA to evaluate the access policy
func getPolicyDecision(policyEngineURL string, path string, input *authzPolicyInput, w io.Writer) (*authzPolicyResult, error) {

	url := policyEngineURL + path

	body, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling json %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf("unexpected content-type: %v", contentType)
	}

	var result authzPolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Result, nil
}

// authorize uses the PAM-configured fields and calls out to OPA to make authorization decisions
func authorize(w io.Writer, uid int, username, url, policyPath, identityFilePath string) AuthResult {

	req := &authzPolicyInput{}
	req.Input.User = username

	// indentifyFilePath is expected to be json, unmarshal it
	raw, err := ioutil.ReadFile(identityFilePath)
	if err != nil {
		fmt.Fprintf(w, "Error reading HostFile: %s\n", err)
		return AuthError
	}
	err = json.Unmarshal(raw, &req.Input.HostIdentity)
	if err != nil {
		fmt.Fprintf(w, "Error decoding HostFile into JSON: %s\n", err)
		return AuthError
	}

	response, err := getPolicyDecision(url, policyPath, req, w)
	if err != nil {
		fmt.Fprintf(w, "Error communicating with the authorization server: %s\n", err)
	}

	if len(response.Errors) == 0 {
		return AuthSuccess
	}

	if len(response.Errors) > 0 {
		fmt.Fprintf(w, "%s \n", response.Errors)
	}

	return AuthError
}

func pamAuthorize(w io.Writer, uid int, username string, argv []string) AuthResult {

	runtime.GOMAXPROCS(1)

	// read the options configured in PAM

	url := defaultUrl
	policyPath := ""
	identityFilePath := defaultHostFile

	for _, arg := range argv {
		opt := strings.Split(arg, "=")
		switch opt[0] {
		case "url":
			url = opt[1]
			pamLog("url set to %s", url)
		case "policy_path":
			policyPath = opt[1]
			pamLog("policy_path set to %s", policyPath)
		case "identity_file_path":
			identityFilePath = opt[1]
			pamLog("identity_file_path set to %s", identityFilePath)
		default:
			pamLog("unkown option: %s\n", opt[0])
		}
	}

	// Call out to OPA for the allow/deny decision

	if url != "" && policyPath != "" {
		return authorize(w, uid, username, url, policyPath, identityFilePath)
	}

	return AuthError
}

func main() {}
