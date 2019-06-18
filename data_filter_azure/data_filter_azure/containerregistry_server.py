#!/usr/bin/env python
import os
import requests
import base64
import json
from flask import Flask,redirect, jsonify, abort, make_response, g, Response
import config
from flask_bootstrap import Bootstrap
import subprocess
import azure.common
from data_filter_azure import opa
import azure.cosmos.cosmos_client as cosmos_client

app = Flask(__name__)
Bootstrap(app)

def check_access_opa(registry_id, user_id, resource_type, resource_name, action):
    decision = opa_get(registry_id, user_id, resource_type, resource_name, action)
    return decision

def opa_get(registry_id, user_id, type, resourceName, action):
    input = {
        'registry': registry_id,
        'user': user_id,
        'type': type,
        'resourceName': resourceName,
        'action': action
    }
    return opa.get_http(path='data/containerregistry/example/allow',
                       input=input)

def oras_pull():
    bundle_name = config.ACR_LOGINSERVER + '/' + config.BUNDLE_TAR + ':' + config.ACR_BUNDLE_ARTIFACT_SHA
    args = ['oras', 'pull', '-u', config.ACR_USERNAME, '-p', config.ACR_PASSWORD, bundle_name]
    output = subprocess.check_output(args, stderr=subprocess.STDOUT)
    print output
    return output


@app.route('/api/registries/<registry_id>/users/<user_id>/<type>/<resource_name>/<action>', methods=["GET"])
def api_check_access(registry_id, user_id, type, resource_name, action):
    return jsonify(check_access_opa(registry_id, user_id, type, resource_name, action))

@app.route('/api/bundles/'+config.BUNDLE_TAR, methods=["GET"])
def api_get_bundle():
    oras_pull()
    with open(os.path.join(app.root_path, '../'+ config.BUNDLE_TAR)) as bundle:
        data=bundle.read()
    return Response(data, content_type="application/gzip")


@app.route('/')
def index():
    return redirect('https://azure.microsoft.com/en-us/services/container-registry/', code = 302)

if __name__ == '__main__':
    app.jinja_env.auto_reload = True
    app.config['TEMPLATES_AUTO_RELOAD'] = True
    app.run(debug=True)
