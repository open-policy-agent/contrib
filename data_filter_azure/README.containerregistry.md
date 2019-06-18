# Data Filter Example with Data from Azure table storage

Usually, OPA keeps policies and data in-memory. But in cases where the context data cannot be loaded, there are sevaral options to support. This directory contains a sample where the data and policies are downloaded as bundles using the [OPA's bundle](https://www.openpolicyagent.org/docs/latest/bundles/) feature. The polcies and data are loaded on the fly and enforced by OPA. Currently OPA is configured to download one bundle from a single source. These bundles are stored in Azure container registry which is OCI-compatible registry. You can read more about this in  [the post](https://stevelasker.blog/2019/01/25/cloud-native-artifact-stores-evolve-from-container-registries/). This sample contains a sample server that uses OPA's Get Document API on the bundle that is configured with it. The server also has an API that downloads the bundle from container registry. OPA is configured to download bundle from this HTTP server and enforce policy using it.

The server itself is implemented in Python using Flask. 

## Running the sample

The is only prototype and doesn't include the authentication. The installation and the model is adopted from the existing sample under data_filter_example of this project

## Data and Policy

A resource called "registry" can contain different types of resources with different supported actions against them. Each authenticated user will have the mappings of allowed resources and actions. For example the below document represents the mappings for a registry owned by user 'bob'. The data represents that Bob can read, write test-repo and read/delete test-chart. 

```json
{
    "bob": {
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

## Running the sample

1. Create the bundle with data and polcies. Create the folder structure with corresponding data and policy rego files.

```bash
export BUNDLE_TAR=registries.tar.gz
cd opabundle
tar -czvf $BUNDLE_TAR .
```

2. Push the bundle to the Azure container registry. Create an [azure container registry](https://azure.microsoft.com/en-us/services/container-registry/)

3. The sample fetchs the container registry keys and login server from environment variables ACR_LOGINSERVER, ACR_USERNAME and ACR_PASSWORD.

```bash
export ACR_LOGINSERVER=<login server of azure container registry>
export ACR_USERNAME=<username of azure container registry>
export ACR_PASSWORD=<password1 of azure container registry>
```

4. Install [oras](https://github.com/deislabs/oras) to push the bundle to the container registry and set the environment variable ACR_BUNDLE_ARTIFACT_SHA to the Digest of the bundle. This will be used to pull the bundle from container registry when OPA makes a request for the bundle. 

```bash
output=$(oras push -u $ACR_USERNAME -p $ACR_PASSWORD $ACR_LOGINSERVER/$BUNDLE_TAR $BUNDLE_TAR)
export ACR_BUNDLE_ARTIFACT_SHA=$(echo $output | grep -Po 'Digest: \K(.*)')
```

5. Install the dependencies into a virtualenv:

```bash
cd ..
virtualenv env
source env/bin/activate
pip install -r requirements.txt
pip install -e .
```

6. Start the server:

```
source env/bin/activate
python data_filter_azure/containerregistry_server.py
```

7. Open a new window and run OPA:
```bash
export BUNDLE_TAR=registries.tar.gz
opa run -s -c opa_config.yaml
```

The server listens on `:5000` and serves an index page by default. It also serves the bundle requests from opa.

8. Send authorization requests to the server which responds with the decision to allow or deny

```bash
curl http://127.0.0.1:5000/api/registries/registry1/users/bob/repositories/repo1/read
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
