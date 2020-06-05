package mongo

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Mongo         *mongo.Client
	Database      string
	ClientOptions *options.ClientOptions
	Logger        *zap.Logger
}

type DBClient interface {
	QueryMongo(pipeline []Queries) (Result, error)
	GetAllDataFromMongo() []EmployeeDetails
}

// Result contains output of partially evaluating a query.
type Result struct {
	Defined bool
	Data    []EmployeeDetails
}

type Queries struct {
	Pipeline bson.M
}

// Create a connection with mongo DB.
func (c *Mongo) CreateConnection() (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), c.ClientOptions)
	if err != nil {
		c.Logger.Error("Unable to establish connection with mongo client",
			zap.String("uri", c.ClientOptions.GetURI()), zap.Error(err))
		return client, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		c.Logger.Error("Unable to ping the mongo db ", zap.Error(err))
		return client, err
	}
	return client, nil
}

// Queries the mongo DB
func (c *Mongo) QueryMongo(pipeline []Queries) (Result, error) {
	var result Result
	var queries []bson.M
	for _, q := range pipeline {
		queries = append(queries, q.Pipeline)
	}
	finalQuery := bson.M{"$or": queries}
	c.Logger.Info("mongo query: ", zap.Any("", finalQuery))
	collection := c.Mongo.Database(c.Database).Collection("employees")
	f, err := collection.Find(context.TODO(), finalQuery)
	if err != nil {
		c.Logger.Error("Failed to query mongo database", zap.Error(err))
		return result, err
	}

	for f.Next(context.Background()) {
		e1 := &EmployeeDetails{}
		err = f.Decode(&e1)
		if err != nil {
			// Do not return on error as we can still try decode other records in range queries.
			c.Logger.Error("Failed to decode the row queried from db.", zap.Error(err))
		}
		result.Data = append(result.Data, *e1)
	}

	result.Defined = true
	return result, nil
}

// Fetches all the employees data from mongo DB.
func (c *Mongo) GetAllDataFromMongo() []EmployeeDetails {
	var records []EmployeeDetails
	f, err := c.Mongo.Database(c.Database).Collection("employees").Find(context.Background(), bson.M{})
	if err != nil {
		c.Logger.Error("Failed to query mongo database", zap.Error(err))
	}

	for f.Next(context.Background()) {
		e1 := &EmployeeDetails{}
		err = f.Decode(&e1)
		if err != nil {
			c.Logger.Error("Failed to decode the row queried from db.", zap.Error(err))
		}
		records = append(records, *e1)
	}

	return records
}

// Parse the == into equivalent mongo query.
func HandleEquals(pipeline *[]bson.M, fieldName string, fieldValue interface{}) {
	filter := bson.M{fieldName: bson.M{"$eq": fieldValue}}
	*pipeline = append(*pipeline, filter)
}

// Parse the != into equivalent mongo query.
func HandleNotEquals(pipeline *[]bson.M, fieldName string, fieldValue interface{}) {
	filter := bson.M{fieldName: bson.M{"$ne": fieldValue}}
	*pipeline = append(*pipeline, filter)
}

// Parse the < into equivalent mongo query.
func HandleLessThan(pipeline *[]bson.M, fieldName string, fieldValue interface{}) {
	filter := bson.M{fieldName: bson.M{"$lt": fieldValue}}
	*pipeline = append(*pipeline, filter)
}

// Parse the > into equivalent mongo query.
func HandleGreaterThan(pipeline *[]bson.M, fieldName string, fieldValue interface{}) {
	filter := bson.M{fieldName: bson.M{"$gt": fieldValue}}
	*pipeline = append(*pipeline, filter)
}

// Parse the <= into equivalent mongo query.
func HandleLessThanEquals(pipeline *[]bson.M, fieldName string, fieldValue interface{}) {
	filter := bson.M{fieldName: bson.M{"$lte": fieldValue}}
	*pipeline = append(*pipeline, filter)
}

// Parse the >= into equivalent mongo query.
func HandleGreaterThanEquals(pipeline *[]bson.M, fieldName string, fieldValue interface{}) {
	filter := bson.M{fieldName: bson.M{"$gte": fieldValue}}
	*pipeline = append(*pipeline, filter)
}
