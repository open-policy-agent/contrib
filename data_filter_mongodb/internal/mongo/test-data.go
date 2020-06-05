package mongo

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type EmployeeDetails struct {
	Name        string
	Designation string
	Salary      int64
	Email       string
	Mobile      string
	Manager     string
}

var (
	e = []EmployeeDetails{
		{Name: "john", Designation: "lead engineer", Salary: 270000, Email: "john@opa.com", Mobile: "1233743438738", Manager: "danerys"},
		{Name: "arya", Designation: "software engineer", Salary: 90000, Email: "arya@opa.com", Mobile: "1233746238738", Manager: "john"},
		{Name: "tyrian", Designation: "senior software engineer", Salary: 250000, Email: "tyrian@opa.com", Mobile: "123336238738", Manager: "danerys"},
		{Name: "jamie", Designation: "lead engineer", Salary: 70000, Email: "jamie@opa.com", Mobile: "1233746238738", Manager: "danerys"},
		{Name: "jeffrey", Designation: "software engineer", Salary: 60000, Email: "jeffrey@opa.com", Mobile: "1233746238738", Manager: "jamie"},
		{Name: "sansa", Designation: "senior software engineer", Salary: 80000, Email: "sansa@opa.com", Mobile: "1233746238738", Manager: "john"},
		{Name: "ramsay", Designation: "software engineer", Salary: 70000, Email: "ramsay@opa.com", Mobile: "1233746238738", Manager: "john"},
		{Name: "cersei", Designation: "senior software engineer", Salary: 170000, Email: "cersei@opa.com", Mobile: "1233746238738", Manager: "jamie"},
		{Name: "theon", Designation: "software engineer", Salary: 75000, Email: "theon@opa.com", Mobile: "1233746238738", Manager: "john"},
		{Name: "rob", Designation: "software engineer", Salary: 75000, Email: "rob@opa.com", Mobile: "1343238738", Manager: "john "},
		{Name: "danerys", Designation: "director of engineering", Salary: 350000, Email: "danerys@opa.com", Mobile: "12332423738"},
	}
)

func (c *Mongo) CreateTestData() {
	collection := c.Mongo.Database(c.Database).Collection("employees")
	records := make([]interface{}, len(e))
	for i := range e {
		records[i] = e[i]
	}
	fmt.Println("inserting the test data...")

	_, err := collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{primitive.E{Key: "name", Value: bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		log.Fatalf("failed to create table %s", err)
	}

	// ignore the error as restarting the opa-mongo server will try to re-insert
	// the data & raise the error.
	_, _ = collection.InsertMany(context.TODO(), records)
}
