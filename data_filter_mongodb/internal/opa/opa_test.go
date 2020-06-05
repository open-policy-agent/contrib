package opa

import (
	"context"
	"log"
	"reflect"
	"testing"

	"github.com/open-policy-agent/contrib/data_filter_mongodb/internal/mongo"
	"go.uber.org/zap"
)

func Test_opaClient_Compile(t *testing.T) {
	type fields struct {
		logger      *zap.Logger
		mongoClient mongo.DBClient
		policyFile  string
	}
	type args struct {
		ctx   context.Context
		input map[string]interface{}
	}

	var f fields
	var err error
	f.logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create zap logger %v", err)
	}
	f.mongoClient = mongo.NewMockMongoClient(nil, "")
	f.policyFile = "./../../example.rego"
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    mongo.Result
		wantErr bool
	}{
		{
			name:   "correct_data_with_unknowns",
			fields: f,
			args: args{
				ctx:   context.Background(),
				input: map[string]interface{}{"method": "GET", "path": []string{"employees", "john"}, "user": "danerys"},
			},
			want: mongo.Result{
				Defined: true,
				Data: []mongo.EmployeeDetails{
					{
						Name:        "john",
						Designation: "Software Engineer",
						Salary:      70000,
						Email:       "john@opa.com",
						Mobile:      "7436238746",
						Manager:     "danerys",
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "incorrect_data_with_unknowns",
			fields: f,
			args: args{
				ctx:   context.Background(),
				input: map[string]interface{}{"method": "GET", "path": []string{"employees", "john"}, "user": "jamie"},
			},
			want: mongo.Result{
				Defined: true,
				Data:    nil,
			},
			wantErr: false,
		},
		{
			name:   "all_empty_no_unknowns",
			fields: f,
			args: args{
				ctx:   context.Background(),
				input: map[string]interface{}{"method": "", "path": []string{"", ""}, "user": ""},
			},
			want: mongo.Result{
				Defined: false,
				Data:    nil,
			},
			wantErr: false,
		},
		{
			name:   "all_nil_no_unknowns",
			fields: f,
			args: args{
				ctx:   context.Background(),
				input: nil,
			},
			want: mongo.Result{
				Defined: false,
				Data:    nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &opaClient{
				logger:      tt.fields.logger,
				mongoClient: tt.fields.mongoClient,
				policyFile:  tt.fields.policyFile,
			}
			got, err := c.Compile(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Compile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
