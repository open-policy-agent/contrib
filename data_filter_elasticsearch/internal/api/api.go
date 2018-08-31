// Copyright 2018 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/olivere/elastic"
	"github.com/open-policy-agent/contrib/data_filter_elasticsearch/internal/es"
	"github.com/open-policy-agent/contrib/data_filter_elasticsearch/internal/opa"
)

const (
	apiCodeNotFound      = "not_found"
	apiCodeParseError    = "parse_error"
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

type Api struct {
	router *mux.Router
	es     *elastic.Client
	index  string
}

func New(esClient *elastic.Client, index string) *Api {

	api := &Api{es: esClient, index: index}
	api.router = mux.NewRouter()

	api.router.HandleFunc("/posts", api.handlGetPosts).Methods(http.MethodGet)
	api.router.HandleFunc("/posts/{id}", api.handleGetPost).Methods(http.MethodGet)

	return api
}

func (api *Api) Run(ctx context.Context) error {
	fmt.Println("Starting server 8080....")
	return http.ListenAndServe(":8080", api.router)
}

func (api *Api) handlGetPosts(w http.ResponseWriter, r *http.Request) {
	result, err := queryOPA(w, r)
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

func (api *Api) handleGetPost(w http.ResponseWriter, r *http.Request) {
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
	combinedQuery := combineQuery(es.GenerateTermQuery("id", vars["id"]), result.Query)
	queryEs(r.Context(), api.es, api.index, combinedQuery, w)
}

func queryOPA(w http.ResponseWriter, r *http.Request) (opa.Result, error) {

	user := r.Header.Get("Authorization")
	path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	input := map[string]interface{}{
		"method": r.Method,
		"path":   path,
		"user":   user,
	}

	return opa.Compile(r.Context(), input)
}

func combineQuery(queryFromHandler elastic.Query, queryFromOpa elastic.Query) elastic.Query {
	var combinedQuery elastic.Query = queryFromHandler
	if queryFromOpa != nil {
		queries := []elastic.Query{queryFromOpa, queryFromHandler}
		combinedQuery = es.GenerateBoolFilterQuery(queries)
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
		Result: es.GetPrettyResult(searchResult),
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
