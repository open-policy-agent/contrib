# Data Filter Example

This directory contains a sample server that uses OPA's Compile API to perform
data filtering and authorization. When the server receives API requests it asks
OPA for a set of conditions to apply to the SQL query that serves the request.

The server itself is implemented in Python using Flask asnd and sqlite3.

## Install

Install the dependencies into a virtualenv:

```bash
virtualenv env
source env/bin/activate
pip install -r requirements.txt
pip install -e .
```

## Testing

Open a new window and run OPA:

```bash
opa run -s example.rego
```

Start the server:

```
source env/bin/activate
python data_filter_example/server.py
```

The server listens on `:5000` and serves an index page by default.

## Development

To run the integration tests, start OPA in another window (`opa run -s`) and
then:

```bash
pytest .
```
