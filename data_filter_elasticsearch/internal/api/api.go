// Copyright 2018 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aquasecurity/esquery"
	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/open-policy-agent/opa/logging"
	"github.com/open-policy-agent/opa/sdk"

	"github.com/gorilla/mux"
	"github.com/open-policy-agent/contrib/data_filter_elasticsearch/internal/es"
	"github.com/open-policy-agent/contrib/data_filter_elasticsearch/internal/opa"
)

const (
	apiCodeInternalError = "internal_error"
	apiCodeNotAuthorized = "not_authorized"
)

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
	opa    *sdk.OPA
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
	fmt.Println("Loading OPA SDK....")
	config, err := os.ReadFile("opa-conf.yaml")
	if err != nil {
		log.Fatal(err)
	}

	logger := logging.New()
	logger.SetLevel(logging.Info)

	opa, err := sdk.New(ctx, sdk.Options{Config: bytes.NewReader(config), Logger: logger})
	if err != nil {
		log.Fatal(err)
	}
	api.opa = opa

	fmt.Println("Starting server 8080....")
	return http.ListenAndServe(":8080", api.router)
}

func (api *ServerAPI) handlGetPosts(w http.ResponseWriter, r *http.Request) {
	result, err := api.queryOPA(w, r)
	if err != nil {
		writeError(w, http.StatusInternalServerError, apiCodeInternalError, err)
		return
	}

	if !result.Defined {
		writeError(w, http.StatusForbidden, apiCodeNotAuthorized, nil)
		return
	}

	combinedQuery := combineQuery(es.GenerateMatchAllQuery(), result.Query)
	queryEs(r.Context(), api.es, api.index, combinedQuery, w)

}

func (api *ServerAPI) handleGetPost(w http.ResponseWriter, r *http.Request) {
	result, err := api.queryOPA(w, r)
	if err != nil {
		writeError(w, http.StatusInternalServerError, apiCodeInternalError, err)
		return
	}

	if !result.Defined {
		writeError(w, http.StatusForbidden, apiCodeNotAuthorized, nil)
		return
	}

	vars := mux.Vars(r)
	combinedQuery := combineQuery(es.GenerateTermQuery("id", vars["id"]), result.Query)
	queryEs(r.Context(), api.es, api.index, combinedQuery, w)
}

func (api *ServerAPI) queryOPA(w http.ResponseWriter, r *http.Request) (opa.Result, error) {

	user := r.Header.Get("Authorization")
	path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	input := map[string]interface{}{
		"method": r.Method,
		"path":   path,
		"user":   user,
	}

	return opa.Compile(r.Context(), api.opa, input)
}

func combineQuery(queryFromHandler esquery.Mappable, queryFromOpa esquery.Mappable) esquery.Mappable {
	var combinedQuery = queryFromHandler
	if queryFromOpa != nil {
		queries := []esquery.Mappable{queryFromOpa, queryFromHandler}
		combinedQuery = es.GenerateBoolFilterQuery(queries)
	}
	return combinedQuery
}

func queryEs(ctx context.Context, client *elastic.Client, index string, query esquery.Mappable, w http.ResponseWriter) {
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
