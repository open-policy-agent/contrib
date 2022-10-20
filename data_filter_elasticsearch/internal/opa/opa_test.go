// Copyright 2018 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package opa

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/open-policy-agent/opa/sdk"
	sdktest "github.com/open-policy-agent/opa/sdk/test"
	"reflect"
	"strings"
	"testing"
)

func TestCompileRequestDeniedAlways(t *testing.T) {
	input := map[string]interface{}{
		"method": "POST",
		"path":   []string{"posts"},
		"user":   "bob",
	}

	policy := `
		package example
		allow = true {
   	 		input.method = "GET"
    		input.path = ["posts"]
		}
	`

	server := sdktest.MustNewServer(
		sdktest.MockBundle("/bundles/bundle.tar.gz", map[string]string{
			"rules.rego": policy,
		}),
	)
	defer server.Stop()

	config := opaConfig(server)

	opa, err := sdk.New(context.Background(), sdk.Options{
		Config: strings.NewReader(config),
	})

	expected := Result{Defined: false}
	result, err := Compile(opa, context.Background(), input)

	if err != nil {
		t.Fatalf("Unexpected error while compiling query: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %v but got: %v", expected, result)
	}
}

func TestCompileRequestAllowedAlways(t *testing.T) {
	input := map[string]interface{}{
		"method": "GET",
		"path":   []string{"posts"},
		"user":   "bob",
	}

	policy := `
		package example
		allow = true {
   	 		input.method = "GET"
    		input.path = ["posts"]
		}
	`

	server := sdktest.MustNewServer(
		sdktest.MockBundle("/bundles/bundle.tar.gz", map[string]string{
			"rules.rego": policy,
		}),
	)
	defer server.Stop()

	config := opaConfig(server)

	opa, err := sdk.New(context.Background(), sdk.Options{
		Config: strings.NewReader(config),
	})

	expected := Result{Defined: true}
	result, err := Compile(opa, context.Background(), input)

	if err != nil {
		t.Fatalf("Unexpected error while compiling query: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %v but got: %v", expected, result)
	}
}

func TestCompileTermQuery(t *testing.T) {
	input := map[string]interface{}{
		"method": "GET",
		"path":   []string{"posts"},
		"user":   "bob",
	}

	policy := `
		package example
		allow = true {
   	 		input.method = "GET"
    		input.path = ["posts"]
			allowed[x]
		}
		
		allowed[x] {
    		x := data.elastic.posts[_]
    		x.author == input.user
		}		
	`

	server := sdktest.MustNewServer(
		sdktest.MockBundle("/bundles/bundle.tar.gz", map[string]string{
			"rules.rego": policy,
		}),
	)
	defer server.Stop()

	config := opaConfig(server)

	opa, err := sdk.New(context.Background(), sdk.Options{
		Config: strings.NewReader(config),
	})

	result, err := Compile(opa, context.Background(), input)

	if err != nil {
		t.Fatalf("Unexpected error while compiling query: %v", err)
	}

	if !result.Defined {
		t.Fatal("Expected result to be defined")
	}

	expectedQueryResult := `{"bool":{"should":[{"term":{"author":{"value":"bob"}}}]}}`

	actualQuerySource := result.Query.Map()
	actualQueryResult, err := marshalQuery(actualQuerySource)
	if err != nil {
		t.Fatalf("Unexpected error while marshalling query: %v", err)
	}

	if actualQueryResult != expectedQueryResult {
		t.Fatalf("Expected %v but got: %v", expectedQueryResult, actualQueryResult)
	}
}

func TestCompileRangeQuery(t *testing.T) {
	input := map[string]interface{}{
		"method":    "GET",
		"path":      []string{"posts"},
		"clearance": 9,
	}

	policy := `
		package example
		allow = true {
   	 		input.method = "GET"
    		input.path = ["posts"]
			allowed[x]
		}
		
		allowed[x] {
    		x := data.elastic.posts[_]
    		x.clearance > input.clearance
		}		
	`

	server := sdktest.MustNewServer(
		sdktest.MockBundle("/bundles/bundle.tar.gz", map[string]string{
			"rules.rego": policy,
		}),
	)
	defer server.Stop()

	config := opaConfig(server)

	opa, err := sdk.New(context.Background(), sdk.Options{
		Config: strings.NewReader(config),
	})

	result, err := Compile(opa, context.Background(), input)

	if err != nil {
		t.Fatalf("Unexpected error while compiling query: %v", err)
	}

	if !result.Defined {
		t.Fatal("Expected result to be defined")
	}

	expectedQueryResult := `{"bool":{"should":[{"range":{"clearance":{"gt":9}}}]}}`

	actualQuerySource := result.Query.Map()
	if err != nil {
		t.Fatalf("Unexpected error while creating query source %v", err)
	}
	actualQueryResult, err := marshalQuery(actualQuerySource)
	if err != nil {
		t.Fatalf("Unexpected error while marshalling query: %v", err)
	}

	if actualQueryResult != expectedQueryResult {
		t.Fatalf("Expected %v but got: %v", expectedQueryResult, actualQueryResult)
	}
}

func TestCompileMustNotQuery(t *testing.T) {
	input := map[string]interface{}{
		"method":    "GET",
		"path":      []string{"posts"},
		"clearance": 9,
	}

	policy := `
		package example
		allow = true {
   	 		input.method = "GET"
    		input.path = ["posts"]
			allowed[x]
		}
		
		allowed[x] {
    		x := data.elastic.posts[_]
    		x.clearance != input.clearance
		}		
	`

	server := sdktest.MustNewServer(
		sdktest.MockBundle("/bundles/bundle.tar.gz", map[string]string{
			"rules.rego": policy,
		}),
	)
	defer server.Stop()

	config := opaConfig(server)

	opa, err := sdk.New(context.Background(), sdk.Options{
		Config: strings.NewReader(config),
	})

	result, err := Compile(opa, context.Background(), input)

	if err != nil {
		t.Fatalf("Unexpected error while compiling query: %v", err)
	}

	if !result.Defined {
		t.Fatal("Expected result to be defined")
	}

	expectedQueryResult := `{"bool":{"should":[{"bool":{"must_not":[{"term":{"clearance":{"value":9}}}]}}]}}`

	actualQuerySource := result.Query.Map()
	actualQueryResult, err := marshalQuery(actualQuerySource)
	if err != nil {
		t.Fatalf("Unexpected error while marshalling query: %v", err)
	}

	if actualQueryResult != expectedQueryResult {
		t.Fatalf("Expected %v but got: %v", expectedQueryResult, actualQueryResult)
	}
}

func TestCompileQueryStringQuery(t *testing.T) {
	input := map[string]interface{}{
		"method":  "GET",
		"path":    []string{"posts"},
		"message": "OPA Rules !",
	}

	policy := `
		package example
		allow = true {
   	 		input.method = "GET"
    		input.path = ["posts"]
			allowed[x]
		}
		
		allowed[x] {
    		x := data.elastic.posts[_]
    		contains(x.message, "OPA")
		}		
	`

	server := sdktest.MustNewServer(
		sdktest.MockBundle("/bundles/bundle.tar.gz", map[string]string{
			"rules.rego": policy,
		}),
	)
	defer server.Stop()

	config := opaConfig(server)

	opa, err := sdk.New(context.Background(), sdk.Options{
		Config: strings.NewReader(config),
	})

	result, err := Compile(opa, context.Background(), input)

	if err != nil {
		t.Fatalf("Unexpected error while compiling query: %v", err)
	}

	if !result.Defined {
		t.Fatal("Expected result to be defined")
	}

	expectedQueryResult := `{"bool":{"should":[{"query_string":{"default_field":"message","query":"*OPA*"}}]}}`

	actualQuerySource := result.Query.Map()
	actualQueryResult, err := marshalQuery(actualQuerySource)
	if err != nil {
		t.Fatalf("Unexpected error while marshalling query: %v", err)
	}

	if actualQueryResult != expectedQueryResult {
		t.Fatalf("Expected %v but got: %v", expectedQueryResult, actualQueryResult)
	}
}

func TestCompileRegexpQuery(t *testing.T) {
	input := map[string]interface{}{
		"method": "GET",
		"path":   []string{"posts"},
		"email":  "j@opa.org",
	}

	policy := `
		package example
		allow = true {
   	 		input.method = "GET"
    		input.path = ["posts"]
			allowed[x]
		}
		
		allowed[x] {
    		x := data.elastic.posts[_]
    		re_match("[a-zA-Z]+@[a-zA-Z]+.org", x.email)
		}		
	`

	server := sdktest.MustNewServer(
		sdktest.MockBundle("/bundles/bundle.tar.gz", map[string]string{
			"rules.rego": policy,
		}),
	)
	defer server.Stop()

	config := opaConfig(server)

	opa, err := sdk.New(context.Background(), sdk.Options{
		Config: strings.NewReader(config),
	})

	result, err := Compile(opa, context.Background(), input)

	if err != nil {
		t.Fatalf("Unexpected error while compiling query: %v", err)
	}

	if !result.Defined {
		t.Fatal("Expected result to be defined")
	}

	expectedQueryResult := `{"bool":{"should":[{"regexp":{"email":{"value":"[a-zA-Z]+@[a-zA-Z]+.org"}}}]}}`

	actualQuerySource := result.Query.Map()
	actualQueryResult, err := marshalQuery(actualQuerySource)
	if err != nil {
		t.Fatalf("Unexpected error while marshalling query: %v", err)
	}

	if actualQueryResult != expectedQueryResult {
		t.Fatalf("Expected %v but got: %v", expectedQueryResult, actualQueryResult)
	}
}

func TestCompileBoolFilterQuery(t *testing.T) {
	input := map[string]interface{}{
		"method":    "GET",
		"path":      []string{"posts"},
		"user":      "bob",
		"clearance": 9,
	}

	policy := `
		package example
		allow = true {
   	 		input.method = "GET"
    		input.path = ["posts"]
			allowed[x]
		}
		
		allowed[x] {
    		x := data.elastic.posts[_]
    		x.author == input.user
			x.clearance > input.clearance
		}		
	`

	server := sdktest.MustNewServer(
		sdktest.MockBundle("/bundles/bundle.tar.gz", map[string]string{
			"rules.rego": policy,
		}),
	)
	defer server.Stop()

	config := opaConfig(server)

	opa, err := sdk.New(context.Background(), sdk.Options{
		Config: strings.NewReader(config),
	})

	result, err := Compile(opa, context.Background(), input)

	if err != nil {
		t.Fatalf("Unexpected error while compiling query: %v", err)
	}

	if !result.Defined {
		t.Fatal("Expected result to be defined")
	}

	expectedQueryResult := `{"bool":{"should":[{"bool":{"filter":[{"term":{"author":{"value":"bob"}}},{"range":{"clearance":{"gt":9}}}]}}]}}`

	actualQuerySource := result.Query.Map()
	actualQueryResult, err := marshalQuery(actualQuerySource)
	if err != nil {
		t.Fatalf("Unexpected error while marshalling query: %v", err)
	}

	if actualQueryResult != expectedQueryResult {
		t.Fatalf("Expected %v but got: %v", expectedQueryResult, actualQueryResult)
	}
}

func TestCompileBoolShouldQuery(t *testing.T) {
	input := map[string]interface{}{
		"method":    "GET",
		"path":      []string{"posts"},
		"user":      "bob",
		"clearance": 9,
	}

	policy := `
		package example
		allow = true {
   	 		input.method = "GET"
    		input.path = ["posts"]
			allowed[x]
		}
		
		allowed[x] {
    		x := data.elastic.posts[_]
    		x.author == input.user
		}

		allowed[x] {
    		x := data.elastic.posts[_]
			x.clearance > input.clearance
		}		
	`

	server := sdktest.MustNewServer(
		sdktest.MockBundle("/bundles/bundle.tar.gz", map[string]string{
			"rules.rego": policy,
		}),
	)
	defer server.Stop()

	config := opaConfig(server)

	opa, err := sdk.New(context.Background(), sdk.Options{
		Config: strings.NewReader(config),
	})

	result, err := Compile(opa, context.Background(), input)

	if err != nil {
		t.Fatalf("Unexpected error while compiling query: %v", err)
	}

	if !result.Defined {
		t.Fatal("Expected result to be defined")
	}

	expectedQueryResult := `{"bool":{"should":[{"term":{"author":{"value":"bob"}}},{"range":{"clearance":{"gt":9}}}]}}`

	actualQuerySource := result.Query.Map()
	actualQueryResult, err := marshalQuery(actualQuerySource)
	if err != nil {
		t.Fatalf("Unexpected error while marshalling query: %v", err)
	}

	if actualQueryResult != expectedQueryResult {
		t.Fatalf("Expected %v but got: %v", expectedQueryResult, actualQueryResult)
	}
}

func marshalQuery(x interface{}) (string, error) {
	d, err := json.Marshal(x)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func opaConfig(server *sdktest.Server) string {
	return fmt.Sprintf(`{
		"services": {
			"test": {
				"url": %q
			}
		},
		"bundles": {
			"test": {
				"resource": "/bundles/bundle.tar.gz"
			}
		},
		"decision_logs": {
			"console": true
		}
	}`, server.URL())
}
