#!/usr/bin/env python

import base64
import os

from flask import Flask
from flask import request
import json
import requests

import logging
import sys

logging.basicConfig(stream=sys.stderr, level=logging.DEBUG)

app = Flask(__name__)

opa_url = os.environ.get("OPA_ADDR", "http://localhost:8181")
policy_path = os.environ.get("POLICY_PATH", "/v1/data/httpapi/authz")


def check_auth(url, user, method, url_as_array, token):
    input_dict = {"input": {"user": user, "path": url_as_array, "method": method}}
    if token is not None:
        input_dict["input"]["token"] = token

    logging.info("Checking auth...")
    logging.info(json.dumps(input_dict, indent=2))
    try:
        rsp = requests.post(url, data=json.dumps(input_dict))
    except Exception as err:
        logging.info(err)
        return {}
    if rsp.status_code >= 300:
        logging.info(
            "Error checking auth, got status %s and message: %s"
            % (j.status_code, j.text)
        )
        return {}
    j = rsp.json()
    logging.info("Auth response:")
    logging.info(json.dumps(j, indent=2))
    return j


@app.route("/", defaults={"path": ""})
@app.route("/<path:path>")
def root(path):
    user_encoded = request.headers.get("Authorization", "Anonymous:none")
    if user_encoded:
        user_encoded = user_encoded.split("Basic ")[1]
    user, _ = base64.b64decode(user_encoded).decode("utf-8").split(":")
    url = opa_url + policy_path
    path_as_array = path.split("/")
    token = request.args["token"] if "token" in request.args else None
    j = check_auth(url, user, request.method, path_as_array, token).get("result", {})
    if j.get("allow", False) == True:
        return "Success: user %s is authorized \n" % user
    return "Error: user %s is not authorized to %s url /%s \n" % (
        user,
        request.method,
        path,
    )


if __name__ == "__main__":
    app.run()
