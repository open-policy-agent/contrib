package mapper

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/open-policy-agent/opa/sdk"

	loggingtest "github.com/open-policy-agent/opa/logging/test"
	sdktest "github.com/open-policy-agent/opa/sdk/test"
)

func TestMapper(t *testing.T) {

	ctx := context.Background()

	server := sdktest.MustNewServer(
		sdktest.MockBundle("/bundles/bundle.tar.gz", map[string]string{
			"main.rego": `
				package test
				allow {
					x := data.elastic.posts[_]
    				x.author == input.y
				}
			`,
		}),
	)

	defer server.Stop()

	config := fmt.Sprintf(`{
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

	testLogger := loggingtest.New()
	opa, err := sdk.New(ctx, sdk.Options{
		Config:        strings.NewReader(config),
		ConsoleLogger: testLogger,
	})
	if err != nil {
		t.Fatal(err)
	}

	defer opa.Stop(ctx)

	var result *sdk.PartialResult
	if result, err = opa.Partial(ctx, sdk.PartialOptions{
		// Path:     "test",
		Input:    map[string]int{"y": 2},
		Query:    "data.test.allow = true",
		Unknowns: []string{"data.elastic"},
		Mapper:   &ElasticMapper{},
		Now:      time.Unix(0, 1619868194450288000).UTC(),
	}); err != nil {
		t.Fatal(err)
	} else if decision, ok := result.Result.(Result); !ok || !decision.Defined {
		t.Fatalf("expected defined result but got %v", decision.Defined)
	}

	entries := testLogger.Entries()

	if l := len(entries); l != 1 {
		t.Fatalf("expected %v but got %v", 1, l)
	}

	// just checking for existence, since it's a complex value
	if entries[0].Fields["mapped_result"] == nil {
		t.Fatalf("expected not nil value for mapped_result but got nil")
	}

	if entries[0].Fields["result"] == nil {
		t.Fatalf("expected not nil value for result but got nil")
	}

	if entries[0].Fields["timestamp"] != "2021-05-01T11:23:14.450288Z" {
		t.Fatalf("expected %v but got %v", "2021-05-01T11:23:14.450288Z", entries[0].Fields["timestamp"])
	}

}
