// Package docker implements helpers to interact with a Docker registry.
package docker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

// Registry represents a Docker registry API to connect to.
type Registry struct {
	baseURL  string
	token    string
	user     string
	password string
}

// Manifest represents a Docker registry image manifest.
type Manifest struct {
	Layers []Layer     `json:"layers,omitempty"`
	Config ImageConfig `json:"config"`
}

// Layer represents a Docker registry layer metadata.
type Layer struct {
	Size   int    `json:"size"`
	Digest string `json:"digest"`
}

// ImageConfig represents a Docker registry image metadata.
type ImageConfig struct {
	Size   int    `json:"size"`
	Digest string `json:"digest"`
}

// New returns a new Registry object.
func New(url, user, password string) *Registry {
	return &Registry{
		baseURL:  url,
		user:     user,
		password: password,
	}
}

// Authorization returns the bearer token to use for registry authentication.
func (r *Registry) Authorization() string {
	if r.token == "" {
		return ""
	}
	return fmt.Sprintf("Bearer %v", r.token)
}

// Path returns the raw layer URL.
func (r *Registry) Path(repo string, layer Layer) string {
	return r.repoURL(repo) + "/blobs/" + layer.Digest
}

// Tags returns a slice of tags for the named repository.
func (r *Registry) Tags(repo string) ([]string, error) {

	req, err := http.NewRequest("GET", r.repoURL(repo)+"/tags/list", nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var body struct {
			Name string   `json:"name"`
			Tags []string `json:"tags,omitempty"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			return nil, err
		}
		return body.Tags, nil
	}

	return nil, fmt.Errorf("unexpected status code: %v", resp.Status)
}

// Manifest returns image manifest information for the repo and tag.
func (r *Registry) Manifest(repo, tag string) (*Manifest, error) {

	url := fmt.Sprintf("%v/manifests/%v", r.repoURL(repo), tag)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	resp, err := r.do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result Manifest

	if resp.StatusCode == http.StatusOK {
		return &result, json.NewDecoder(resp.Body).Decode(&result)
	}

	return nil, fmt.Errorf("unexpected status code: %v", err)
}
func (r *Registry) repoURL(repo string) string {
	return fmt.Sprintf("%v/v2/%v", r.baseURL, repo)
}

func (r *Registry) login(unauthorized *http.Response) error {
	h := unauthorized.Header.Get("Www-Authenticate")
	parts := token.FindStringSubmatch(h)
	realm, service, scope := parts[1], parts[2], parts[3]
	url := fmt.Sprintf("%v?service=%v&scope=%v", realm, service, scope)
	if r.user != "" {
		url += "&account=" + r.user
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	if r.user != "" {
		req.SetBasicAuth(r.user, r.password)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from login: %v", resp.Status)
	}
	var body struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return err
	}
	r.token = body.Token
	return nil
}

func (r *Registry) do(req *http.Request) (*http.Response, error) {
	return r.doImpl(req, true)
}

func (r *Registry) doImpl(req *http.Request, authenticate bool) (*http.Response, error) {

	authz := r.Authorization()
	if authz != "" {
		req.Header.Add("Authorization", authz)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		if authenticate {
			if err := r.login(resp); err != nil {
				return nil, err
			}
			return r.doImpl(req, false)
		}
	}

	return resp, nil
}

// Example:
//
// Bearer realm="https://auth.docker.io/token",service="registry.docker.io",scope="repository:openpolicyagent/opa:pull"
var token = regexp.MustCompile(`Bearer realm="(.*?)",service="(.*?)",scope="(.*?)"`)
