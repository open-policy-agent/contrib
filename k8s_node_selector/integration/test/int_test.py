import unittest
import random
import warnings
import time
import requests
import subprocess
from kubernetes import client, config
from kubernetes.client.apis import core_v1_api
from kubernetes.client.rest import ApiException


class TestStringMethods(unittest.TestCase):
    def test_opa_api_secured(self):
        warnings.filterwarnings("ignore", category=ResourceWarning)

        # Forward port to the OPA instance in KIND
        port_forward_command = "kubectl port-forward service/opa -n opa 8181:443"
        port_forward_process = subprocess.Popen(port_forward_command.split())

        # Wait for port forward to start
        time.sleep(3)

        try:
            # Check /health is allowed for anonymous traffic
            response = requests.get("https://localhost:8181/health",  verify=False)
            self.assertEqual(response.status_code, 200, "Expect health endpoint to return 200")


            # Check access to policies
            response = requests.get("https://localhost:8181/v1/policies", verify=False)
            self.assertEqual(response.status_code, 401, "Expect authz policy to refuse access to policies api")
        finally:
            port_forward_process.kill()
        

    def test_create_pod_check_node_selector_set(self):
        warnings.filterwarnings("ignore", category=ResourceWarning)

        config.load_kube_config()
        core_v1 = core_v1_api.CoreV1Api()
        mapped_namespace = self.create_mapped_namespace(core_v1)
        # time.sleep(30)
        spec = {
            "containers": [
                {
                    "image": "busybox",
                    "name": "sleep",
                    "args": ["/bin/sh", "-c", "while true;do date;sleep 5; done"],
                }
            ]
        }
        podInCluster = self.createPodWithSpec(core_v1, spec, mapped_namespace)
        self.assertIsNotNone(podInCluster.spec.node_selector)
        self.assertEqual(podInCluster.spec.node_selector["agentpool"], "pool1")

    def test_create_pod_in_ignored_namespace_and_no_node_selector_applied(self):
        warnings.filterwarnings("ignore", category=ResourceWarning)

        config.load_kube_config()
        core_v1 = core_v1_api.CoreV1Api()
        spec = {
            "containers": [
                {
                    "image": "busybox",
                    "name": "sleep",
                    "args": ["/bin/sh", "-c", "while true;do date;sleep 5; done"],
                }
            ]
        }
        podInCluster = self.createPodWithSpec(core_v1, spec, namespace="kube-system")
        self.assertIsNone(podInCluster.spec.node_selector)

    def test_pod_with_node_selector_already_expect_fail(self):
        warnings.filterwarnings("ignore", category=ResourceWarning)

        config.load_kube_config()
        core_v1 = core_v1_api.CoreV1Api()

        mapped_namespace = self.create_mapped_namespace(core_v1)

        spec = {
            "nodeSelector": {"agentpool": "notallowedvalue"},
            "containers": [
                {
                    "image": "busybox",
                    "name": "sleep",
                    "args": ["/bin/sh", "-c", "while true;do date;sleep 5; done"],
                }
            ],
        }

        with self.assertRaises(ApiException) as context:
            self.createPodWithSpec(core_v1, spec, mapped_namespace)

        self.assertEqual(context.exception.reason, "Forbidden")
        self.assertTrue(
            "Manually specifying NodeSelector not supported in this namespace."
            in context.exception.body,
            f"Didn't find expected error message in response. Actual response:{context.exception.body}",
        )

    def test_pod_in_unmapped_namespace_expect_allow(self):
        warnings.filterwarnings("ignore", category=ResourceWarning)

        config.load_kube_config()
        core_v1 = core_v1_api.CoreV1Api()

        unmapped_namespace = "unmapped" + str(random.randint(0, 1000))
        core_v1.create_namespace(
            client.V1Namespace(metadata=client.V1ObjectMeta(name=unmapped_namespace))
        )

        # Wait for the new namespace to be usable
        self.wait_for_namespace_active(core_v1, unmapped_namespace)
        self.wait_for_service_account_ready(core_v1, unmapped_namespace)

        spec = {
            "containers": [
                {
                    "image": "busybox",
                    "name": "sleep",
                    "args": ["/bin/sh", "-c", "while true;do date;sleep 5; done"],
                }
            ]
        }

        podInCluster = self.createPodWithSpec(core_v1, spec, unmapped_namespace)
        self.assertIsNone(podInCluster.spec.node_selector)

    def create_mapped_namespace(self, core_v1):
        mapped_namespace = "mapped" + str(random.randint(0, 1000))
        core_v1.create_namespace(
            client.V1Namespace(
                metadata=client.V1ObjectMeta(
                    name=mapped_namespace, labels={"agentpool": "pool1"}
                )
            )
        )

        # Wait for the new namespace to be usable
        self.wait_for_namespace_active(core_v1, mapped_namespace)
        self.wait_for_service_account_ready(core_v1, mapped_namespace)
        return mapped_namespace

    def createPodWithSpec(self, core_v1, spec, namespace="default"):
        pod_name = "testpod" + str(random.randint(0, 1000))

        pod_manifest = {
            "apiVersion": "v1",
            "kind": "Pod",
            "metadata": {"name": pod_name, "namespace": namespace},
            "spec": spec,
        }

        return core_v1.create_namespaced_pod(body=pod_manifest, namespace=namespace)

    def wait_for_namespace_active(self, core_v1, unmapped_namespace):
        # Wait for the new namespace to be active
        namespace = core_v1.read_namespace_status(unmapped_namespace)
        while namespace.status.phase != "Active":
            print("Waiting for namespace to be ready")
            time.sleep(1)
            namespace = core_v1.read_namespace_status(unmapped_namespace)

    def wait_for_service_account_ready(self, core_v1, unmapped_namespace):
        # Wait for the default service account to be valid in the namespace
        service_accounts = core_v1.list_namespaced_service_account(unmapped_namespace)
        while (
            len(service_accounts.items) < 1 or service_accounts.items[0].secrets is None
        ):
            print("Waiting for default service account to be ready")
            time.sleep(1)
            service_accounts = core_v1.list_namespaced_service_account(
                unmapped_namespace
            )


if __name__ == "__main__":
    unittest.main(verbosity=2)
