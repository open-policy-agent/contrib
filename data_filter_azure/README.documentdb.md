# Data Filter Example With Rego-SQL on Azure CosmosDB 

This directory contains a sample server that uses OPA's Compile API to perform
data filtering and authorization. When the server receives API requests it asks
OPA for a set of conditions to apply to the SQL query. The SQl query is applied on the Azure Cosmos Document DB and based on the results, the server responds with the boolean value of whether to allow or deny the request.

The server itself is implemented in Python using Flask. 

## Data and Policy

A resource called "registry" can contain different types of resources with different supported actions against them. Each authenticated user will have a registry with the mappings of allowed resources and actions. For example the below document represents the mappings for a registry owned by user 'bob'. The map represents that Bob can read, write repo1 and read/write chart1 and some wild card permissions.

```json
{
    "registry" : "registry1",
    "user": "bob",
    "d": "blob",
    "map": [
        {
            "type": "repositories",
            "name": "repo1",
            "actions": ["read", "write"]
        },
        {
            "type": "repositories",
            "name": "repo2",
            "actions": ["*"]
        },
        {
            "type": "charts",
            "name": "chart1",
            "actions": ["read", "write"]
        },
        {
            "type": "pipelines",
            "name": "*",
            "actions": ["read"]
        }
    ]
}
```

An azure Cosmos DB can have JSON documents and the Rego-SQL be directly applied on the JSON documents using the Document DB APIs.

## Running the sample

The is only prototype and doesn't include the authentication. The installation and the model is adopted from the existing sample under data_filter_example of this project

1. Create an azure cosmos db (SQL API) and create a azure table [Azure Cosmos DB](https://docs.microsoft.com/en-us/azure/cosmos-db/introduction)

2. Install the dependencies into a virtualenv:

```bash
virtualenv env
source env/bin/activate
pip install -r requirements.txt
pip install -e .
```

3. The sample fetches the COSMOS DB endpoint and primary key from environment variables COSMOSDB_ENDPOINT and COSMOSDB_PRIMARYKEY
 
```bash
export COSMOSDB_ENDPOINT=<cosmos db endpoint>
export COSMOSDB_PRIMARYKEY=<cosmos db primary key>
```

4. Open a new window and run OPA:
```bash
opa run -s example_documentdb.rego
```

5. Start the server:

```
source env/bin/activate
python data_filter_azure/documentdb_server.py
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
