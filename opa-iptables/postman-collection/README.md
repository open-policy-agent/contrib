# OPA API postman collection

Postman collection demonstrating how to use [OPA](https://www.openpolicyagent.org/) REST API. More information on API can be found at [https://www.openpolicyagent.org/docs/latest/rest-api/](https://www.openpolicyagent.org/docs/latest/rest-api/)

## Installation

To use the latest publish version of a collection, click the following button to import the OPA REST API as a collection:

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/a37391cb5b1d88e10a58)

You can also download collection from this repo, and manually import `.json` collection file via `File > import` in postman.

### Prerequisites

- *Postman* The collection is for use by the Postman app. Postman is a utility that allows you to quickly test and use REST APIs. More information can be found at [getpostman.com](https://www.getpostman.com/).
- You need OPA server running locally for experimenting this REST APIs. Use [Docker](https://www.docker.com/) for running OPA server. Following is a command to run OPA with docker:
    >`$ docker run -p 8181:8181 openpolicyagent/opa:0.10.7 run --server --log-level=debug`

    > If you don't have docker install then follow the instruction from [here](https://docs.docker.com/install/) to download docker for your platform.

## See Also

[OPA API Documentaion](https://www.openpolicyagent.org/docs/latest/rest-api/)

[Postman API development tool](https://www.getpostman.com/)