# OPA Data Filtering on Azure

This directory contains sample servers that use OPA for data filtering and authorization against Azure storage services. The servers are implemented in Python using Flask.

Two variants are provided, each targeting a different Azure storage backend:

- **[Cosmos DB (SQL API)](./README.documentdb.md)** - Uses OPA's Compile API to generate SQL conditions that are applied directly on Cosmos DB queries
- **[Table Storage](./README.tablestorage.md)** - Uses OPA's Get Document API with input to query authorization data from Azure Table Storage at evaluation time

See the linked READMEs above for setup instructions and usage details for each variant.
