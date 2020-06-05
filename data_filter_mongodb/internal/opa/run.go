package opa

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	mongo2 "github.com/open-policy-agent/contrib/data_filter_mongodb/internal/mongo"
	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

// ServerAPI is the Server's API.
type ServerAPI struct {
	Router    *mux.Router
	OpaClient *opaClient
}

// New return the server's API.
func New(mongoClient *mongo2.Mongo, p string) *ServerAPI {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create zap logeer %v", err)
	}
	api := &ServerAPI{
		OpaClient: &opaClient{
			mongoClient: mongoClient,
			logger:      logger,
			policyFile:  p,
		},
	}
	api.Router = mux.NewRouter()
	api.Router.HandleFunc("/employees", api.handleGetReq).Methods(http.MethodGet)
	api.Router.HandleFunc("/employees/{id}", api.handleGetReq).Methods(http.MethodGet)
	api.Router.HandleFunc("/records", api.getAllRecords)

	return api
}

// Run the server.
func (c *ServerAPI) Run(ctx context.Context) error {
	fmt.Println("Starting server 9095....")
	return http.ListenAndServe(":9095", c.Router)
}

//
//func (c *ServerAPI) handleGetReqs(w http.ResponseWriter, r *http.Request) {
//	result, err := c.OpaClient.queryMongoOPA(r)
//	if err != nil {
//		_ = fmt.Errorf("failed to handle the request %s", err)
//		return
//	}
//
//	if !result.Defined {
//		return
//	}
//	writeJSON(w, http.StatusOK, apiWrapper{
//		Result: result,
//	})
//}

// Get all records from mongo DB.
func (c *ServerAPI) getAllRecords(w http.ResponseWriter, r *http.Request) {
	records := c.OpaClient.mongoClient.GetAllDataFromMongo()
	resultBytes, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write(resultBytes)
	if err != nil {
		_ = fmt.Errorf("failed to write received result to response writer %s", err)
	}
}

type apiWrapper struct {
	Result interface{} `json:"result"`
}

func (c *ServerAPI) handleGetReq(w http.ResponseWriter, r *http.Request) {
	result, err := c.OpaClient.queryMongoOPA(r)
	if err != nil {
		_ = fmt.Errorf("failed to handle the request %s", err)
		return
	}

	if !result.Defined {
		return
	}
	writeJSON(w, http.StatusOK, apiWrapper{
		Result: result,
	})
}

func writeJSON(w http.ResponseWriter, status int, x interface{}) {
	bs, _ := json.Marshal(x)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write(bs)
	if err != nil {
		_ = fmt.Errorf("failed to write json result set to response writer %s", err)
	}
}
