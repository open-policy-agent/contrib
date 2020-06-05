package mongo

import (
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MockMongoClientStruct struct {
	Mongo         *mongo.Client
	Database      string
	ClientOptions *options.ClientOptions
	Logger        *zap.Logger
}

func NewMockMongoClient(mongo *mongo.Client, database string) MockMongoClient {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create zap logger %v", err)
	}
	return &MockMongoClientStruct{
		Mongo:    mongo,
		Database: database,
		Logger:   logger,
	}
}

func (c *MockMongoClientStruct) GetAllDataFromMongo() []EmployeeDetails {
	return e
}

type MockMongoClient interface {
	QueryMongo(pipeline []Queries) (Result, error)
	GetAllDataFromMongo() []EmployeeDetails
}

const (
	// data for successful case
	johnDanerys = "[{map[$and:[map[name:map[$eq:danerys]] map[name:map[$eq:john]]]]} {map[$and:[map[manager:map[$eq:danerys]] map[name:map[$eq:john]]]]}]"
	// data for unsuccessful case
	johnJamie = "[{map[$and:[map[name:map[$eq:jamie]] map[name:map[$eq:john]]]]} {map[$and:[map[manager:map[$eq:jamie]] map[name:map[$eq:john]]]]}]"
)

func (c *MockMongoClientStruct) QueryMongo(pipeline []Queries) (Result, error) {
	var r Result
	if johnDanerys == fmt.Sprintf("%s", pipeline) {
		r = Result{
			Defined: true,
			Data: []EmployeeDetails{
				{
					Name:        "john",
					Designation: "Software Engineer",
					Salary:      70000,
					Email:       "john@opa.com",
					Mobile:      "7436238746",
					Manager:     "danerys",
				},
			},
		}
	} else if johnJamie == fmt.Sprintf("%s", pipeline) {
		r = Result{
			Defined: true,
			Data:    nil,
		}
	}
	return r, nil
}
