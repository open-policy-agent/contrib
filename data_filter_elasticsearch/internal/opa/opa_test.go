// Copyright 2018 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package opa

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
)

func TestCompileRequestDeniedAlways(t *testing.T) {
	input := map[string]interface{}{
		"method": "GET",
		"path":   []string{"post"},
		"user":   "bob",
	}

	policy := `
		package example
		allow = true {
   	 		input.method = "GET"
    		input.path = ["posts"]
		}
	`

	expected := Result{Defined: false}
	result, err := Compile(context.Background(), input, []byte(policy))

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

	expected := Result{Defined: true}
	result, err := Compile(context.Background(), input, []byte(policy))

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

	result, err := Compile(context.Background(), input, []byte(policy))

	if err != nil {
		t.Fatalf("Unexpected error while compiling query: %v", err)
	}

	if !result.Defined {
		t.Fatal("Expected result to be defined")
	}

	expectedQueryResult := `{"bool":{"_name":"BoolShouldQuery","should":{"term":{"author":{"_name":"TermQuery","value":"bob"}}}}}`

	actualQuerySource, err := result.Query.Source()
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

	result, err := Compile(context.Background(), input, []byte(policy))

	if err != nil {
		t.Fatalf("Unexpected error while compiling query: %v", err)
	}

	if !result.Defined {
		t.Fatal("Expected result to be defined")
	}

	expectedQueryResult := `{"bool":{"_name":"BoolShouldQuery","should":{"range":{"clearance":{"from":9,"include_lower":false,"include_upper":true,"to":null}}}}}`

	actualQuerySource, err := result.Query.Source()
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

	result, err := Compile(context.Background(), input, []byte(policy))

	if err != nil {
		t.Fatalf("Unexpected error while compiling query: %v", err)
	}

	if !result.Defined {
		t.Fatal("Expected result to be defined")
	}

	expectedQueryResult := `{"bool":{"_name":"BoolShouldQuery","should":{"bool":{"_name":"BoolMustNotQuery","must_not":{"term":{"clearance":9}}}}}}`

	actualQuerySource, err := result.Query.Source()
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

	result, err := Compile(context.Background(), input, []byte(policy))

	if err != nil {
		t.Fatalf("Unexpected error while compiling query: %v", err)
	}

	if !result.Defined {
		t.Fatal("Expected result to be defined")
	}

	expectedQueryResult := `{"bool":{"_name":"BoolShouldQuery","should":{"query_string":{"_name":"QueryStringQuery","default_field":"message","query":"*OPA*"}}}}`

	actualQuerySource, err := result.Query.Source()
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

	result, err := Compile(context.Background(), input, []byte(policy))

	if err != nil {
		t.Fatalf("Unexpected error while compiling query: %v", err)
	}

	if !result.Defined {
		t.Fatal("Expected result to be defined")
	}

	expectedQueryResult := `{"bool":{"_name":"BoolShouldQuery","should":{"regexp":{"email":{"value":"[a-zA-Z]+@[a-zA-Z]+.org"}}}}}`

	actualQuerySource, err := result.Query.Source()
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

	result, err := Compile(context.Background(), input, []byte(policy))

	if err != nil {
		t.Fatalf("Unexpected error while compiling query: %v", err)
	}

	if !result.Defined {
		t.Fatal("Expected result to be defined")
	}

	expectedQueryResult := `{"bool":{"_name":"BoolShouldQuery","should":{"bool":{"_name":"BoolFilterQuery","filter":[{"term":{"author":{"_name":"TermQuery","value":"bob"}}},{"range":{"clearance":{"from":9,"include_lower":false,"include_upper":true,"to":null}}}]}}}}`

	actualQuerySource, err := result.Query.Source()
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

	result, err := Compile(context.Background(), input, []byte(policy))

	if err != nil {
		t.Fatalf("Unexpected error while compiling query: %v", err)
	}

	if !result.Defined {
		t.Fatal("Expected result to be defined")
	}

	expectedQueryResult := `{"bool":{"_name":"BoolShouldQuery","should":[{"term":{"author":{"_name":"TermQuery","value":"bob"}}},{"range":{"clearance":{"from":9,"include_lower":false,"include_upper":true,"to":null}}}]}}`

	actualQuerySource, err := result.Query.Source()
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

func marshalQuery(x interface{}) (string, error) {
	d, err := json.Marshal(x)
	if err != nil {
		return "", err
	}
	return string(d), nil
}
