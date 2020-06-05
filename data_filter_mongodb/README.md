# OPA-MongoDB Data filtering

This is a cli to integrate OPA with mongo database.

### How it works:

The motive behind writing the OPA mongo DB data filtering service is to have an 
example implementation of how we can leverage OPA in enforcing policies 
on data retrieval with mongo database. This is an application which queries 
mongo database whereas users or other applications query this service using
a http request. This service needs a config file with authentication details
for mongo DB and a rego policy file which contains the policies we want to 
enforce on mongo DB. Every request made to this service
is evaluated with rego policy and the request is parsed based on the `user input` + 
`rego policy` in to equivalent `mongo query` to retrieve the desired data. If the 
request doesn't comply with rego policy the service sends back an empty 
response.

This service helps to enforce policies on mongo database without changing the 
business logic as per the requirements.

Available commands & flags:
```
$ ./opa-mongo --help
This CLI tool integrates OPA with MongoDB. It leverages OPA's partial evaluation feature to translate Rego to a MongoDB query which can then be applied to the database.

Note: Before you run this CLI make sure you have working mongo database that is accessible through this cli.

Usage:
  opa_mongo [command]

Available Commands:
  help        Help about any command
  run         Run OPA with Mongo DB as server

Flags:
  -c, --config.file string   configuration file.
  -h, --help                 help for opa_mongo

Use "opa-mongo [command] --help" for more information about a command.
```

You can deploy opa-mongo server in two ways:

1. Docker
2. Kubernetes

### Deploying as container:

Prerequisites:

* Docker

#### Step 1:

Running Mongo DB as container:

```
docker run -d -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=password -p 27017:27017 mongo
```

#### Step 2:

Configure the mongo db details in the config file as shown below in a file named ```config.yaml```:

**Note:** In step 1 mongo DB container is created with root user as `admin` and root password as `password`, so, 
we are providing the same auth credentials in the below config file.

```
cfg:
  address: <node-ip>:27017
  database: opa-mongo
  username: admin
  password: password
``` 

#### Step 3: 

Add example rego policy into the file named ```example.rego```:

```
package details.authz

allow {
  input.method == "GET"
  employee := data.employees[_]
  employee.name == input.user
  employee.name == input.path[1]
}

allow {
  input.method == "GET"
  employee := data.employees[_]
  employee.manager == input.user
  employee.name == input.path[1]
}

allow {
  input.method == "GET"
  input.path = ["employees"]
  employee := data.employees[_]
  input.user == "danerys"
  employee.salary > 0
  employee.salary < 300000
}
```

The above rego policy validates that the requested employee details can be retrieved by the employee and by a respective manager.

#### Step 4:

Below cmd lets you to start opa server with mongo db and on opa server startup, we insert sample data.

```
$ docker run -d -p 9095:9095 -v /config.yaml:/opa/config.yaml /
    -v /example.rego:/opa/example.rego vineeth97/opa-mongo:latest
```



### Deploying in Kubernetes:

Pre-requisites:

1. kubectl
2. kind - https://kind.sigs.k8s.io/

#### Step 1:

Create a Kubernetes cluster 

```
kind create cluster
```

#### Step 2:

Running opa-mongo server with mongo database

```
kubectl create -f https://github.com/open-policy-agent/contrib/blob/master/data_filter_mongodb/k8s/opa-mongo.yaml
```

By now you have running opa-mongo server with mongo database. Now let's query the data from opa-mongo server.

### View the data that exists in mongo database:

You can view all the data that exist in the mongo database by accessing the below URL. This endpoint is exposed for viewing the data that exists in the database.

```
http://localhost:9095/records
```

This opa-mongo server is designed to retrieve employee details. The employee details can be seen either by the employee or by a respective manager. 

Below are records that are inserted on opa server startup:

```
[
  {
    "Name": "john",
    "Designation": "lead engineer",
    "Salary": 270000,
    "Email": "john@opa.com",
    "Mobile": "1233743438738",
    "Manager": "danerys"
  },
  {
    "Name": "arya",
    "Designation": "software engineer",
    "Salary": 90000,
    "Email": "arya@opa.com",
    "Mobile": "1233746238738",
    "Manager": "john"
  },
  {
    "Name": "tyrian",
    "Designation": "senior software engineer",
    "Salary": 250000,
    "Email": "tyrian@opa.com",
    "Mobile": "123336238738",
    "Manager": "danerys"
  },
  {
    "Name": "jamie",
    "Designation": "lead engineer",
    "Salary": 70000,
    "Email": "jamie@opa.com",
    "Mobile": "1233746238738",
    "Manager": "danerys"
  },
  {
    "Name": "jeffrey",
    "Designation": "software engineer",
    "Salary": 60000,
    "Email": "jeffrey@opa.com",
    "Mobile": "1233746238738",
    "Manager": "jamie"
  },
  {
    "Name": "sansa",
    "Designation": "senior software engineer",
    "Salary": 80000,
    "Email": "sansa@opa.com",
    "Mobile": "1233746238738",
    "Manager": "john"
  },
  {
    "Name": "ramsay",
    "Designation": "software engineer",
    "Salary": 70000,
    "Email": "ramsay@opa.com",
    "Mobile": "1233746238738",
    "Manager": "john"
  },
  {
    "Name": "cersei",
    "Designation": "senior software engineer",
    "Salary": 170000,
    "Email": "cersei@opa.com",
    "Mobile": "1233746238738",
    "Manager": "jamie"
  },
  {
    "Name": "theon",
    "Designation": "software engineer",
    "Salary": 75000,
    "Email": "theon@opa.com",
    "Mobile": "1233746238738",
    "Manager": "john"
  },
  {
    "Name": "rob",
    "Designation": "software engineer",
    "Salary": 75000,
    "Email": "rob@opa.com",
    "Mobile": "1343238738",
    "Manager": "john"
  },
  {
    "Name": "danerys",
    "Designation": "director of engineering",
    "Salary": 350000,
    "Email": "danerys@opa.com",
    "Mobile": "12332423738",
  }
]
```

### Querying opa-mongo server:

Below is the request you can make to opa server:

This requests for employee john details as we notice request param and path in request body says john, and the request body expects the requested user name. Based on the rego policy the requested user should be either the employee or the respective manager to retrieve the employee details.

```
URL: http://localhost:9095/employees/john
```
Request body:
```
{
    "input": {
        "method": "GET",
        "path": ["employees", "john"],
        "user": "danerys"
    }
}
```

Response:
```
{
    "result": {
        "Defined": true,
        "Data": [
            {
                "Name": "john",
                "Designation": "lead engineer",
                "Salary": 270000,
                "Email": "john@opa.com",
                "Mobile": "1233743438738",
                "Manager": "danerys"
            }
        ]
    }
}
```

In opa-mongo server logs you can observe the received query, equivalent opa-query and the equivalent mongo-query based on the rego policy:

```
2020/08/07 18:15:45 received requestmap[method:GET path:[employees john] user:danery]
2020/08/07 18:15:45 opa-query: ["danery" = data.employees[_].name; "john" = data.employees[_].name "danery" = data.employees[_].manager; "john" = data.employees[_].name]
2020/08/07 18:15:45 mongo query: map[$or:[map[$and:[map[name:map[$eq:danery]] map[name:map[$eq:john]]]] map[$and:[map[manager:map[$eq:danery]] map[name:map[$eq:john]]]]]]
```

Now you have successfully retrieved john details by requesting the opa-mongo server by enforcing the policies from the mongo database.

### Supported OPA Built-in Functions:

#### Comparison:

- [x] ==
- [x] !=
- [x] <
- [x] <=
- [x] >
- [x] >=

## Support for OPA references

References are used to access nested documents in OPA. OPA policies can be written over deeply nested structures which the server would then translate to Mongo `Nested` queries.

## Generated Mongo queries

For the OPA operators mentioned above, following are Mongo queries generated by the server:

### Term level Queries

- Term Query
- Range Query

### Joining Queries

- Nested Query

### Compound Queries

- Bool Query

### Limitations:

- The OPA policies should be written according to the fields in the MongoDB documents to get the desired results.
- The server supports limited OPA operators and returns an empty result if the OPA policy contains an unsupported operator and logs the error.
- The server supports only two endpoints `/employees` and `/employees/{employee_name}` for fetching posts created when the server starts.