#!/usr/bin/env python
import requests
import base64
import json
from flask import Flask,redirect, jsonify, abort, make_response, g
import config
from flask_bootstrap import Bootstrap
import azure.common
from data_filter_azure import opa
import azure.cosmos.cosmos_client as cosmos_client

app = Flask(__name__)
Bootstrap(app)

def check_access_opa(registry_id, user_id, type, resource_name, action):
    decision = query_opa(registry_id, user_id, type, resource_name, action)
    if not decision.defined:
        raise abort(403)

    sql = opa.splice(SELECT='permissions.id', FROM='permissions JOIN map in permissions.map', WHERE=None, decision=decision)
    print(sql)
    result = query_cosmosdb(sql, args=None, one=True)
    if len(result) == 0:
        return False
    return True


@app.route('/api/registries/<registry_id>/users/<user_id>/<type>/<resource_name>/<action>', methods=["GET"])
def api_check_access(registry_id, user_id, type, resource_name, action):
    return jsonify(check_access_opa(registry_id, user_id, type, resource_name, action))


@app.route('/')
def index():
    return redirect('https://docs.microsoft.com/en-us/azure/cosmos-db/introduction', code = 302)

def query_cosmosdb(query, args=[], one=False):
    dbinfo = get_cosmosdb()
    cosmosdbquery = {
                "query": query
            }
    options = {}
    options['enableCrossPartitionQuery'] = True
    options['maxItemCount'] = 2
    client = dbinfo['client']
    container = dbinfo['container']
    result_iterable = client.QueryItems(container['_self'], cosmosdbquery, options)
    values = []
    for item in iter(result_iterable):
        return item
        values.append(item)
    return values

def query_opa(registry_id, user_id, type, resourceName, action):
    input = {
        'registry': registry_id,
        'user': user_id,
        'type': type,
        'resourceName': resourceName,
        'action': action
    }
    return opa.compile(q='data.documentdb.example.allow==true',
                       input=input,
                       unknowns=['permissions'])

def get_cosmosdb():
    dbinfo = dict();
    client = cosmos_client.CosmosClient(url_connection=config.COSMOSDB_ENDPOINT, auth={
                                    'masterKey': config.COSMOSDB_PRIMARYKEY})
    dbinfo['client'] = client
    id = config.COSMOSDB_DATABASE
    databases = list(client.QueryDatabases({
            "query": "SELECT * FROM r WHERE r.id=@id",
            "parameters": [
                { "name":"@id", "value": id }
            ]
        }))

    if len(databases) > 0:
        db = databases[0]
    else:
        db = client.CreateDatabase({'id': id})
    dbinfo['db'] = db
    containerid = 'permissions'
    database_link = 'dbs/' + id
    collections = list(client.QueryContainers(
            database_link,
            {
                "query": "SELECT * FROM r WHERE r.id=@id",
                "parameters": [
                    { "name":"@id", "value": containerid }
                ]
            }
        ))

    if len(collections) > 0:
        container = collections[0]
    else:
        options = {
            'offerThroughput': 400
        }
        container_definition = {
            'id': containerid,
            'partitionKey': {'paths': ['/registry']}
        }
        container = client.CreateContainer(db['_self'], container_definition, options)

    dbinfo['container'] = container
    return dbinfo

def add_documents():
    dbinfo = get_cosmosdb()
    client = dbinfo['client']
    container = dbinfo['container']
    for document in DOCUMENTS:
        client.UpsertItem(container['_self'], document)

def init_db():
    with app.app_context():
        add_documents()

DOCUMENTS = [
    {
        'registry' : 'registry1',
        'user': 'bob',
        'id': 'blob',
        'map': [
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
    },
    {
        'registry' : 'registry1',
        'user': 'alice',
        'id': 'alice',
        'map': [
            {
                "type": "repositories",
                "name": "*",
                "actions": ["*"]
            },
            {
                "type": "charts",
                "name": "chart1",
                "actions": ["read"]
            }
        ]
    }
]

if __name__ == '__main__':
    init_db()
    app.jinja_env.auto_reload = True
    app.config['TEMPLATES_AUTO_RELOAD'] = True
    app.run(debug=True)
