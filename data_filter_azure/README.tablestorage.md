# Data Filter Example with Data from Azure table storage

Usually, OPA keeps policies and data in-memory. But in cases where the context data cannot be loaded, there are sevaral options to support. This directory contains a sample where the data is queried on-fly from Azure table storage and is provided as the input  to OPA when it executes the policy query. This sample contains a sample server that uses OPA's Get Document API. Once the user is authenticated, it queries the authorization data for the user from azure storage and executes the policy queries using it. 
This directory contains a sample server that uses OPA's Get Document API (with input) to perform data filtering and authorization. When the server receives auhtorization requests it first queries the context from Azure table storage. It then includes this data as input to the OPA execute policy query. It then responds with the boolean value of whether to allow or deny the request.

The server itself is implemented in Python using Flask. 

## Data and Policy

A resource called "registry" can contain different types of resources with different supported actions against them. Each authenticated user will have a registry with the mappings of allowed resources and actions. For example the below document represents the mappings for a registry owned by user 'bob'. The map represents that Bob can read, write test-repo and read/delete test-chart

```json
{
    "registry": "registry1",
    "user": "bob",
    "map": {
        "repositories": {
            "repo1": {
                "actions": ["read", "write"]
            },
            "repo2": {
                "actions": ["*"]
            }
        },
        "charts": {
            "chart1": {
                "actions": ["read", "write"]
            }
        },
        "pipelines": {
            "*": {
                "actions": ["read"]
            }
        }
    }
}
```

As azure table storage is indexed by partition key and row key, it is required that these JSON documents are stored with  registry id as partition key and user id as the row key. When the server receives the request, it pulls the registry and user information, use them as partition key and row key and fetch the JSON document from table storage and sends as input to the OPA engine to execute the policy

## Running the sample

The is only prototype and doesn't include the authentication. The installation and the model is adopted from the existing sample under data_filter_example of this project

1. Create an azure storage account and create a azure table [Azure Table Storage](https://docs.microsoft.com/en-us/azure/cosmos-db/table-storage-overview)

2. Install the dependencies into a virtualenv:

```bash
virtualenv env
source env/bin/activate
pip install -r requirements.txt
pip install -e .
```

3. The sample fetches the storage connection string from environment variable "AZURE_STORAGE_CONNECTION_STRING" and the table name is "TABLE_NAME"
 
```bash
export AZURE_STORAGE_CONNECTION_STRING=<actual connection string>
export TABLE_NAME=<table_name>
```

4. Open a new window and run OPA:
```bash
opa run -s example_tablestorage.rego
```

5. Start the server:

```
source env/bin/activate
python data_filter_azure/tablestorage_server.py
```

The server listens on `:5000` and serves an index page by default.

6. Send authorization requests to the server which responds with the decision to allow or deny

```bash
curl http://127.0.0.1:5000/api/registries/r1/users/user1/repos/myrepo/read
true
```

| Request | Decision Output |
| --- | --- |
| /api/registries/registry1/users/bob/repositories/repo1/read | true |
| /api/registries/registry1/users/bob/repositories/repo1/write | true |
| /api/registries/registry1/users/bob/repositories/my-repo/write | false |
| /api/registries/registry1/users/bob/repositories/repo2/write | true |
| /api/registries/registry1/users/bob/charts/chart1/read | true |
| /api/registries/registry1/users/bob/charts/test-chart/read | false |
| /api/registries/registry1/users/bob/pipelines/mypipeline/read | true |
| /api/registries/registry1/users/bob/pipelines/mypipeline/write | false |
