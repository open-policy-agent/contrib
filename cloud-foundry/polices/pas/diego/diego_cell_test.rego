package main_test

import rego.v1

import data.main

test_vm_type_is_present if {
	expected := {
		"resource": [{
			"identifier": "diego_cell",
			"description": "Diego Cell",
			"instances": "",
			"instances_best_fit": 0,
			"instance_type_id": "medium",
			"instance_type_best_fit": "xlarge.disk",
		}],
		"present": true,
	}

	actual := main.find_vm_type(mock_api_call, "diego_cell")

	actual == expected
}

test_vm_type_is_absent if {
	expected := {
		"resource": [],
		"present": false,
	}

	actual := main.find_vm_type(mock_bad_api_call, "diego_cell")

	actual == expected
}

mock_bad_api_call := {"resources": []}

mock_api_call := {"resources": [
	{
		"identifier": "mysql_monitor",
		"description": "Monitors the MySQL Cluster",
		"instances": "",
		"instances_best_fit": 1,
		"instance_type_id": "",
		"instance_type_best_fit": "micro",
	},
	{
		"identifier": "clock_global",
		"description": "Schedules asynchronous tasks for cloud controller",
		"instances": "",
		"instances_best_fit": 0,
		"instance_type_id": "",
		"instance_type_best_fit": "medium.disk",
	},
	{
		"identifier": "cloud_controller_worker",
		"description": "Worker for cloud controller asynchronous tasks",
		"instances": "",
		"instances_best_fit": 0,
		"instance_type_id": "",
		"instance_type_best_fit": "micro",
	},
	{
		"identifier": "diego_cell",
		"description": "Diego Cell",
		"instances": "",
		"instances_best_fit": 0,
		"instance_type_id": "medium",
		"instance_type_best_fit": "xlarge.disk",
	},
]}
