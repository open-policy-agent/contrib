# Kubernetes Local Validation using Styra

## Directions
### Lab Preparation
1. Clone the required repository:
    ```shell
    git clone https://github.com/StyraInc/academy-samples.git
    cd academy-samples/kubernetes/repo
    ```

### Administrator setting up the manifests repo with validation configs
1. Create the repo validation config file.
    File name:
    ```
    .styra.yaml
    ```

    Content:
    ```yaml
    checks:
      dev_deployments:
        system_id: e9136611bd7842cf92d1e508c15e78c4
        files:
            - k8s/deployments/*.yaml
    ```

2. Commit change to the repo.
    ```shell
    git add .styra.yaml
    git commit -m "Styra config to check deployment YAMLs"
    ```

3. (Pretend for training purpose) Merge the change into remote so that this configuration will be checked out with the repo.

4. (Trainer) Add a validating policy in Styra workspace `K8s hands-on` system to require readiness probe.

5. (Trainer) Create an API token for local validation. Grant `SystemViewer` permission to the token.

### Developer adding manifests to the repo
#### Set-up local toolchain for validation

1. Download the styra CLI from the docs page: [https://docs.styra.com/das/reference/cli/install-use-cli](https://docs.styra.com/das/reference/cli/install-use-cli)

2. Confirm the CLI is working:
    ```shell
    styra version
    ```

3. Configure the connection to Styra workpsace by creating the config file in your user directory.
    File path:
    ```
    ~/.styra/config
    ```

    Content:
    ```yaml
    organization_id: edu-das.styra.com
    token: <secret token to be provided by trainer>
    ```
#### Create and validate manifest
1. Create a resource manifest `k8s/deployments/nginx.yaml` to add an nginx ingress deployment. View the file that is already created for you.
    ```shell
	  cat k8s/deployments/nginx.yaml
	```

2. Locally validate new manifest before submission for review.
    ```shell
    styra validate check-local
    ```

    Sample output
    ```console
    % styra validate check-local
    Check dev_deployments failed

    Found errors in the following files:
    * k8s/deployments/nginx.yaml
    ```

    To see more detailed output:
    ```shell
    styra validate check-local --output out.json
    cat out.json
    ```

#### Update and validate manifest
1. Add a readiness probe to the new deployment to pass validation. Uncomment the following section in `k8s/deployments/nginx.yaml`:
    ```yaml
    # readinessProbe:
    #   httpGet:
    #     path: /nginx-ready
    #     port: readiness-port
    #   periodSeconds: 1
    ```

2. Validate the updated manifest.
    ```console
    % styra validate check-local
    Check dev_deployments passed: 1 file checked.
    ```

3. (Pretend for training purpose) Submit the manifest for review.