#!/usr/bin/env python3

"""Converts the output of OPA's JSON test results to JUnit XML and
prints it to stdout.

Usage:
  opa_test_to_junit.py <path>

Where path is a file containing the output of `opa test --format=json ...`

The script also accept input on stdin allowing the output to be piped:
  opa test --format=json <path> | opa_test_to_junit.py
"""

# JUnit format used as documented here: https://llg.cubic.org/docs/junit/

import sys
import json
import socket
import fileinput
import xml.etree.ElementTree as ET

def _nanos_to_seconds(nanos):
    time_ms = nanos / 1000000
    return round(time_ms / 1000, 3)

def main(json_doc):
    """Read provided OPA JSON test report and print it as JUnit XML"""
    test_suites = {}
    json_report = json.loads(json_doc)
    for element in json_report:
        if element['package'] not in test_suites:
            test_suites[element['package']] = {
                'testcases': [],
                'time': 0,
                'tests': 0,
                'failures': 0,
                'skipped': 0,
                'errors': 0
            }

        package = test_suites[element['package']]

        package['tests'] += 1
        package['time'] += element['duration']

        testcase = ET.Element(
            "testcase",
            name=element['name'],
            time=str(_nanos_to_seconds(element['duration'])),
            classname=element['location']['file']
        )

        if 'fail' in element and element['fail'] is True:
            package['failures'] += 1
            failure = ET.Element("failure")
            testcase.append(failure)

        if 'skip' in element:
            package['skipped'] += 1
            skipped = ET.Element("skipped")
            testcase.append(skipped)

        if 'error' in element:
            package['errors'] += 1
            error = ET.Element(
                "error",
                type=element['error']['code'],
                message=element['error']['message']
            )
            testcase.append(error)

        package['testcases'].append(testcase)

    total_metrics = {
        'time': 0,
        'tests': 0,
        'failures': 0,
        'skipped': 0,
        'errors': 0
    }
    el_testsuites = ET.Element("testsuites")

    for suite_name in test_suites:
        suite = test_suites[suite_name]
        el_testsuite = ET.Element(
            "testsuite",
            name=suite_name,
            hostname=socket.gethostname(),
            time=str(_nanos_to_seconds(suite['time'])),
            tests=str(suite['tests']),
            failures=str(suite['failures']),
            skipped=str(suite['skipped']),
            errors=str(suite['errors'])
        )
        for testcase in suite['testcases']:
            el_testsuite.append(testcase)

        total_metrics['tests'] += len(suite['testcases'])
        total_metrics['time'] += suite['time']
        total_metrics['failures'] += suite['failures']
        total_metrics['skipped'] += suite['skipped']
        total_metrics['errors'] += suite['errors']

        el_testsuites.append(el_testsuite)

    el_testsuites.set("time", str(_nanos_to_seconds(total_metrics['time'])))
    el_testsuites.set("tests", str(total_metrics['tests']))
    el_testsuites.set("failures", str(total_metrics['failures']))
    el_testsuites.set("skipped", str(total_metrics['skipped']))
    el_testsuites.set("errors", str(total_metrics['errors']))

    tree = ET.ElementTree(el_testsuites)
    tree.write(sys.stdout.fileno(), encoding='utf-8', method='xml', xml_declaration=True)

if __name__ == '__main__':
    if len(sys.argv) == 2 and (sys.argv[1] == '--help' or sys.argv[1] == 'help'):
        print(__doc__)
        sys.exit(0)

    main(''.join(fileinput.input()))
