package mongo

import (
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func TestMongoClientStruct_GetAllDataFromMongo(t *testing.T) {
	type fields struct {
		Mongo         *mongo.Client
		Database      string
		ClientOptions *options.ClientOptions
		Logger        *zap.Logger
	}
	tests := []struct {
		name   string
		fields fields
		want   []EmployeeDetails
	}{
		{
			name: "all-data",
			want: e,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MockMongoClientStruct{
				Mongo:         tt.fields.Mongo,
				Database:      tt.fields.Database,
				ClientOptions: tt.fields.ClientOptions,
				Logger:        tt.fields.Logger,
			}
			if got := c.GetAllDataFromMongo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllDataFromMongo() = %v, want %v", got, tt.want)
			}
		})
	}
}
