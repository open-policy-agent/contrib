#!/usr/bin/env python3

"""Send CloudFormation Template(s) to OPA for validation"""

# pylint: disable=import-error,invalid-name,anomalous-backslash-in-string,broad-except

import argparse
import json
import os
import sys

import requests

from cfn_flip import to_json

def get_all_templates(directory):
    """Get all template files in directory"""
    templates = []

    for root, _, files in os.walk(directory):
        for file in files:
            if file.endswith((".json", ".yml", ".yaml")):
                templates.append(os.path.join(root, file))

    return templates


def check_resource(path, name, resource, expect_allow, test_mode):
    """Checks a single resource from a template against OPA"""
    properties = bools_to_string(resource["Properties"])

    opa_input = {
        "action": "CREATE",
        "hook": "OpenPolicyAgent::OPA::Hook",
        "resource": {
            "id": name,
            "name": resource["Type"],
            "type": resource["Type"],
            "properties": properties
        }
    }

    resp = requests.post("http://localhost:8181", json=opa_input)

    decision = resp.json()

    if test_mode:
        if decision["allow"] == expect_allow:
            print(f"SUCCESS: {path} {name}")

            return True

        print(f"FAIL: {path} {name}")

        if len(decision["violations"]) > 0:
            print()

        for violation in decision["violations"]:
            print(f"\t{violation}")

        if len(decision["violations"]) > 0:
            print()

        return False

    print(json.dumps(decision, indent=4))

    return decision["allow"]

def check_template(path, test_mode):
    """Send each resource in provided template to OPA for validation"""
    contents = {}
    try:
        with open(path, encoding='UTF-8') as file:
            if path.endswith(".json"):
                contents = json.loads(file.read())
            else:
                contents = json.loads(to_json(file.read(), clean_up=True))

    except Exception as e:
        print(f"ERROR: Exception raised when loading {path}", e)
        return False

    resources = contents["Resources"]
    resource_names = list(resources.keys())

    success = True
    for name in resource_names:
        # Some templates contain resources which we don't intend to test
        # For those we may use an "ObjectToTest" attribute in the templates
        # metadata section to point out which object should be verified.
        if "Metadata" in contents and "ObjectToTest" in contents["Metadata"]:
            if name == contents["Metadata"]["ObjectToTest"]:
                if not check_resource(path, name, resources[name], "success" in path, test_mode):
                    success = False
            else:
                continue

        if not check_resource(path, name, resources[name], "success" in path, test_mode):
            success = False

    return success


def bools_to_string(obj):
    """
    When presented to the hook, AWS templates converts booleans to strings...

    ¯\_(ツ)_/¯

    """
    for k, v in obj.items():
        if isinstance(v, dict):
            bools_to_string(v)
        if isinstance(v, list):
            for item in v:
                if isinstance(item, dict):
                    bools_to_string(item)

        if isinstance(v, bool):
            obj[k] = "false" if not v else "true"

    return obj


def main():
    """Validate on or more templates against OPA running on localhost

    If --test=true is provided, test all templates by asserting file names containing "success"
    are allowed, and others expected to fail.

    Example usage:

    $ validate.py ../examples/templates/eks-cluster-logging/eks-deny-cluster.yaml
    {
        "allow": false,
        "violations": [
            "no logging types are enabled for cluster: EksCluster",
            "public endpoint needs to be disabled for cluster: EksCluster"
        ]
    }

    $ validate.py --test=true ../examples/templates
    SUCCESS: iam-fail-no-user-policy.yaml UserWithNoPolicies
    SUCCESS: iam-success-user-policy-attached.yaml UserWithPolicy
    SUCCESS: iam-success-user-policy-attached.yaml UserWithManagedPolicy
    ...
    """
    parser = argparse.ArgumentParser()

    parser.add_argument("-t", "--test", type=bool, default=False, help="Test for failures")
    parser.add_argument("files", nargs="*")

    args = parser.parse_args()

    if len(args.files) == 0:
        sys.exit("No files provided. Example usage: validate.py file1.yaml file2.json dir1/")

    templates = []
    for file in args.files:
        if os.path.isdir(file):
            templates += get_all_templates(file)
        else:
            templates.append(file)

    success = True
    for template in templates:
        if not check_template(os.path.relpath(template), args.test):
            success = False

    if not success:
        sys.exit(1)

if __name__ == "__main__":
    main()
