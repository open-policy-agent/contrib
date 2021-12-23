// Copyright 2018 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/olivere/elastic"
	"github.com/open-policy-agent/contrib/data_filter_elasticsearch/internal/es"
	"github.com/open-policy-agent/contrib/data_filter_elasticsearch/internal/resolvers"
	"github.com/open-policy-agent/opa/sdk"
)

const (
	apiCodeNotFound      = "not_found"
	apiCodeParseError    = "parse_error"
	apiCodeInternalError = "internal_error"
	apiCodeNotAuthorized = "not_authorized"
)

var opa *sdk.OPA

type apiError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message,omitempty"`
	} `json:"error"`
}

type apiWrapper struct {
	Result interface{} `json:"result"`
}

// ServerAPI is the Server's API.
type ServerAPI struct {
	router *mux.Router
	es     *elastic.Client
	index  string
}

// New return the server's API.
func New(esClient *elastic.Client, index string) *ServerAPI {

	api := &ServerAPI{es: esClient, index: index}
	api.router = mux.NewRouter()

	api.router.HandleFunc("/posts", api.handlGetPosts).Methods(http.MethodGet)
	api.router.HandleFunc("/posts/{id}", api.handleGetPost).Methods(http.MethodGet)

	return api
}

// Run the server.
func (api *ServerAPI) Run(ctx context.Context) error {
	opa = startOpa()
	fmt.Println("Starting server 8080....")
	return http.ListenAndServe(":8080", api.router)
}

func (api *ServerAPI) handlGetPosts(w http.ResponseWriter, r *http.Request) {
	result, err := queryOPA(w, r)
	if err != nil {
		writeError(w, http.StatusInternalServerError, apiCodeInternalError, err)
		return
	}

	if !result.Defined {
		writeError(w, http.StatusForbidden, apiCodeNotAuthorized, nil)
		return
	}

	combinedQuery := combineQuery(resolvers.GenerateMatchAllQuery(), result.Query)
	queryEs(r.Context(), api.es, api.index, combinedQuery, w)

}

func (api *ServerAPI) handleGetPost(w http.ResponseWriter, r *http.Request) {
	result, err := queryOPA(w, r)
	if err != nil {
		writeError(w, http.StatusInternalServerError, apiCodeInternalError, err)
		return
	}

	if !result.Defined {
		writeError(w, http.StatusForbidden, apiCodeNotAuthorized, nil)
		return
	}

	vars := mux.Vars(r)
	combinedQuery := combineQuery(resolvers.GenerateTermQuery("id", vars["id"]), result.Query)
	queryEs(r.Context(), api.es, api.index, combinedQuery, w)
}

func startOpa() *sdk.OPA {
	config, err := os.ReadFile("opa-conf.yaml")
	if err != nil {

		panic(err)
	}
	opa, err := sdk.New(context.Background(), sdk.Options{
		Config: bytes.NewReader(config),
	})
	if err != nil {
		panic(err)
	}
	return opa
}

func queryOPA(w http.ResponseWriter, r *http.Request) (resolvers.Result, error) {

	user := r.Header.Get("Authorization")
	path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	input := map[string]interface{}{
		"method": r.Method,
		"path":   path,
		"user":   user,
	}

	decision, err := opa.Partial(r.Context(), sdk.PartialOptions{
		Input:    input,
		Unknowns: []string{"data.elastic"},
		Path:     "example/allow",
		Query:    "data.example.allow == true",
		Resolver: &resolvers.ElasticResolver{},
	})
	if err != nil {
		return resolvers.Result{}, err
	}
	return decision.Result.(resolvers.Result), nil
}

func combineQuery(queryFromHandler elastic.Query, queryFromOpa elastic.Query) elastic.Query {
	var combinedQuery elastic.Query = queryFromHandler
	if queryFromOpa != nil {
		queries := []elastic.Query{queryFromOpa, queryFromHandler}
		combinedQuery = resolvers.GenerateBoolFilterQuery(queries)
	}
	return combinedQuery
}

func queryEs(ctx context.Context, client *elastic.Client, index string, query elastic.Query, w http.ResponseWriter) {
	searchResult, err := es.ExecuteEsSearch(ctx, client, index, query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, apiCodeInternalError, err)
		return
	}

	writeJSON(w, http.StatusOK, apiWrapper{
		Result: es.GetPrettyESResult(searchResult),
	})
	return
}

func writeError(w http.ResponseWriter, status int, code string, err error) {
	var resp apiError
	resp.Error.Code = code
	if err != nil {
		resp.Error.Message = err.Error()
	}
	writeJSON(w, status, resp)
}

func writeJSON(w http.ResponseWriter, status int, x interface{}) {
	bs, _ := json.Marshal(x)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(bs)
}
