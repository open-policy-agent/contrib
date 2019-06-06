#!/usr/bin/env python
import requests
import base64
import json
from flask import Flask,redirect, jsonify, abort, make_response, g
import config
from flask_bootstrap import Bootstrap
import azure.common
from data_filter_azure import opa
from azure.storage import CloudStorageAccount
from tablestorageaccount import TableStorageAccount
from azure.storage.table import TableService, Entity

app = Flask(__name__)
table_name = config.TABLE_NAME
Bootstrap(app)

def check_access_opa(registry_id, user_id, resource_type, resource_name, action):
    table_service = get_table_service()
    mapEntity = table_service.get_entity(table_name, registry_id, user_id)
    if mapEntity is None:
        raise abort(404)

    decision = opa_get(registry_id, user_id, mapEntity, resource_type, resource_name, action)
    return decision

def opa_get(registry_id, user_id, mapEntity, type, resourceName, action):
    input = {
        'registry': registry_id,
        'user': user_id,
        'type': type,
        'resourceName': resourceName,
        'action': action,
        'data': json.loads(mapEntity['map'])
    }
    return opa.get_http(path='data/example/allow',
                       input=input)

@app.route('/api/registries/<registry_id>/users/<user_id>/<type>/<resource_name>/<action>', methods=["GET"])
def api_check_access(registry_id, user_id, type, resource_name, action):
    return jsonify(check_access_opa(registry_id, user_id, type, resource_name, action))


@app.route('/')
def index():
    return redirect('https://docs.microsoft.com/en-us/azure/cosmos-db/table-storage-overview', code = 302)

def get_table_service():
    table_service = getattr(g, '_table_service', None)
    if table_service is None:
        if config.IS_EMULATED:
            account = TableStorageAccount(is_emulated=True)
        else:
            account_connection_string = config.STORAGE_CONNECTION_STRING
            # Split into key=value pairs removing empties, then split the pairs into a dict
            configuration = dict(s.split('=', 1) for s in account_connection_string.split(';') if s)

            # Authentication
            account_name = configuration.get('AccountName')
            account_key = configuration.get('AccountKey')
            # Basic URL Configuration
            endpoint_suffix = configuration.get('EndpointSuffix')
            if endpoint_suffix == None:
                table_endpoint  = configuration.get('TableEndpoint')
                table_prefix = '.table.'
                start_index = table_endpoint.find(table_prefix)
                end_index = table_endpoint.endswith(':') and len(table_endpoint) or table_endpoint.rfind(':')
                endpoint_suffix = table_endpoint[start_index+len(table_prefix):end_index]
                account = TableStorageAccount(account_name = account_name, connection_string = account_connection_string, endpoint_suffix=endpoint_suffix)

        table_service = g._table_service = account.create_table_service()
    return table_service

def add_table_entities():
    table_service = get_table_service()
    table_service.create_table(table_name)
    for entity in ENTITIES:
        table_entity = {'PartitionKey': entity['registry'], 'RowKey': entity['user'], 'map' : json.dumps(entity['map'])}
        table_service.insert_or_replace_entity(table_name, table_entity)

def init_table():
    with app.app_context():
        add_table_entities()

ENTITIES = [
    {
        'registry' : 'registry1',
        'user': 'bob',
        'map': {
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
    },
    {
        'registry' : 'registry1',
        'user': 'alice',
        'map': {
            "repositories": {
                "*": {
                    "actions": ["*"]
                }
            },
            "charts": {
                "chart1": {
                    "actions": ["read"]
                }
            }
        }
    },
    {
        'registry' : 'registry2',
        'user': 'charlie',
        'map': {
            "repositories": {
                "repo2": {
                    "actions": ["read", "write"]
                }
            },
            "charts": {
                "chart1": {
                    "actions": ["read"]
                }
            },
            "pipelines": {
                "*": {
                    "actions": ["*"]
                }
            }
        }
    },
    {
        'registry' : 'registry2',
        'user': 'sam',
        'map': {
            "charts": {
                "chart1": {
                    "actions": ["read", "write", "delete"]
                }
            },
            "pipelines": {
                "pipeline1": {
                    "actions": ["read"]
                }
            }
        }
    }
]

if __name__ == '__main__':
    init_table()
    app.jinja_env.auto_reload = True
    app.config['TEMPLATES_AUTO_RELOAD'] = True
    app.run(debug=True)
