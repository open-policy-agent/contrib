package node_selector.testdata

import rego.v1

example_namespaces_with_label := {"default": {
	"apiVersion": "v1",
	"kind": "Namespace",
	"metadata": {
		"creationTimestamp": "2020-01-17T14:18:28Z",
		"labels": {"agentpool": "pool1"},
		"name": "default",
		"resourceVersion": "517",
		"selfLink": "/api/v1/namespaces/kube-system",
		"uid": "d20094e6-bdfb-4d7c-a887-a8fd99d9f3dc",
	},
	"spec": {"finalizers": ["kubernetes"]},
	"status": {"phase": "Active"},
}}

example_pod_has_node_selector := {
	"apiVersion": "admission.k8s.io/v1beta1",
	"kind": "AdmissionReview",
	"request": {
		"kind": {
			"group": "",
			"kind": "Pod",
			"version": "v1",
		},
		"namespace": "default",
		"object": {
			"metadata": {
				"creationTimestamp": "2018-10-27T02:12:20Z",
				"labels": {"app": "nginx"},
				"name": "nginx",
				"namespace": "default",
				"uid": "bbfee96d-d98d-11e8-b280-080027868e77",
			},
			"spec": {
				"nodeSelector": {"nodePool": "pool1"},
				"containers": [{
					"image": "nginx",
					"imagePullPolicy": "Always",
					"name": "nginx",
					"resources": {},
					"terminationMessagePath": "/dev/termination-log",
					"terminationMessagePolicy": "File",
					"volumeMounts": [{
						"mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
						"name": "default-token-tm9v8",
						"readOnly": true,
					}],
				}],
				"dnsPolicy": "ClusterFirst",
				"restartPolicy": "Always",
				"schedulerName": "default-scheduler",
				"securityContext": {},
				"serviceAccount": "default",
				"serviceAccountName": "default",
				"terminationGracePeriodSeconds": 30,
				"tolerations": [
					{
						"effect": "NoExecute",
						"key": "node.kubernetes.io/not-ready",
						"operator": "Exists",
						"tolerationSeconds": 300,
					},
					{
						"effect": "NoExecute",
						"key": "node.kubernetes.io/unreachable",
						"operator": "Exists",
						"tolerationSeconds": 300,
					},
				],
				"volumes": [{
					"name": "default-token-tm9v8",
					"secret": {"secretName": "default-token-tm9v8"},
				}],
			},
			"status": {
				"phase": "Pending",
				"qosClass": "BestEffort",
			},
		},
		"oldObject": null,
		"operation": "CREATE",
		"resource": {
			"group": "",
			"resource": "pods",
			"version": "v1",
		},
		"uid": "bbfeef88-d98d-11e8-b280-080027868e77",
		"userInfo": {
			"groups": [
				"system:masters",
				"system:authenticated",
			],
			"username": "minikube-user",
		},
	},
}

example_pod_has_nonexisting_namespace := {
	"apiVersion": "admission.k8s.io/v1beta1",
	"kind": "AdmissionReview",
	"request": {
		"kind": {
			"group": "",
			"kind": "Pod",
			"version": "v1",
		},
		"namespace": "barrywhite",
		"object": {
			"metadata": {
				"creationTimestamp": "2018-10-27T02:12:20Z",
				"labels": {"app": "nginx"},
				"name": "nginx",
				"namespace": "barrywhite",
				"uid": "bbfee96d-d98d-11e8-b280-080027868e77",
			},
			"spec": {
				"containers": [{
					"image": "nginx",
					"imagePullPolicy": "Always",
					"name": "nginx",
					"resources": {},
					"terminationMessagePath": "/dev/termination-log",
					"terminationMessagePolicy": "File",
					"volumeMounts": [{
						"mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
						"name": "default-token-tm9v8",
						"readOnly": true,
					}],
				}],
				"dnsPolicy": "ClusterFirst",
				"restartPolicy": "Always",
				"schedulerName": "default-scheduler",
				"securityContext": {},
				"serviceAccount": "default",
				"serviceAccountName": "default",
				"terminationGracePeriodSeconds": 30,
				"tolerations": [
					{
						"effect": "NoExecute",
						"key": "node.kubernetes.io/not-ready",
						"operator": "Exists",
						"tolerationSeconds": 300,
					},
					{
						"effect": "NoExecute",
						"key": "node.kubernetes.io/unreachable",
						"operator": "Exists",
						"tolerationSeconds": 300,
					},
				],
				"volumes": [{
					"name": "default-token-tm9v8",
					"secret": {"secretName": "default-token-tm9v8"},
				}],
			},
			"status": {
				"phase": "Pending",
				"qosClass": "BestEffort",
			},
		},
		"oldObject": null,
		"operation": "CREATE",
		"resource": {
			"group": "",
			"resource": "pods",
			"version": "v1",
		},
		"uid": "bbfeef88-d98d-11e8-b280-080027868e77",
		"userInfo": {
			"groups": [
				"system:masters",
				"system:authenticated",
			],
			"username": "minikube-user",
		},
	},
}

example_pod_doesnt_have_node_selector := {
	"apiVersion": "admission.k8s.io/v1beta1",
	"kind": "AdmissionReview",
	"request": {
		"kind": {
			"group": "",
			"kind": "Pod",
			"version": "v1",
		},
		"namespace": "default",
		"object": {
			"metadata": {
				"creationTimestamp": "2018-10-27T02:12:20Z",
				"labels": {"app": "nginx"},
				"name": "default",
				"namespace": "default",
				"uid": "bbfee96d-d98d-11e8-b280-080027868e77",
			},
			"spec": {
				"containers": [{
					"image": "nginx",
					"imagePullPolicy": "Always",
					"name": "nginx",
					"resources": {},
					"terminationMessagePath": "/dev/termination-log",
					"terminationMessagePolicy": "File",
					"volumeMounts": [{
						"mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
						"name": "default-token-tm9v8",
						"readOnly": true,
					}],
				}],
				"dnsPolicy": "ClusterFirst",
				"restartPolicy": "Always",
				"schedulerName": "default-scheduler",
				"securityContext": {},
				"serviceAccount": "default",
				"serviceAccountName": "default",
				"terminationGracePeriodSeconds": 30,
				"tolerations": [
					{
						"effect": "NoExecute",
						"key": "node.kubernetes.io/not-ready",
						"operator": "Exists",
						"tolerationSeconds": 300,
					},
					{
						"effect": "NoExecute",
						"key": "node.kubernetes.io/unreachable",
						"operator": "Exists",
						"tolerationSeconds": 300,
					},
				],
				"volumes": [{
					"name": "default-token-tm9v8",
					"secret": {"secretName": "default-token-tm9v8"},
				}],
			},
			"status": {
				"phase": "Pending",
				"qosClass": "BestEffort",
			},
		},
		"oldObject": null,
		"operation": "CREATE",
		"resource": {
			"group": "",
			"resource": "pods",
			"version": "v1",
		},
		"uid": "bbfeef88-d98d-11e8-b280-080027868e77",
		"userInfo": {
			"groups": [
				"system:masters",
				"system:authenticated",
			],
			"username": "minikube-user",
		},
	},
}
