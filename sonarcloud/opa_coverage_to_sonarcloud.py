import argparse
from dataclasses import dataclass, field
from typing import List
import json
import xml.etree.ElementTree as ET
import os
import time
import math


@dataclass
class RowInfo:
    row: int


@dataclass
class LineCoverage:
    start: RowInfo
    end: RowInfo


@dataclass
class BaseCoverage:
    coverage: int
    covered_lines: int = 0
    not_covered_lines: int = 0


@dataclass
class FileCoverage(BaseCoverage):
    covered: List[LineCoverage] = field(default_factory=list)
    not_covered: List[LineCoverage] = field(default_factory=list)


@dataclass
class OPACoverage(BaseCoverage):
    files: dict[str, FileCoverage] = field(default_factory=dict)

    @property
    def lines_valid(self):
        return self.covered_lines + self.not_covered_lines


def generate_sonarcloud_xml(coverage: OPACoverage) -> ET.Element:
    coverage_xml = ET.Element("coverage", attrib={
        "version": "1",
    })
    for path, data in coverage.files.items():
        if isinstance(data, dict):
            if not "coverage" in data:
                data["coverage"] = 0
            data = FileCoverage(**data)
        file_data = ET.Element("file", attrib={
            "path": path,
        })
        coverage_xml.append(file_data)
        cover_map = {}
        max_line = 1
        for c in data.covered:
            if isinstance(c, dict):
                c = LineCoverage(**c)
                if isinstance(c.start, dict):
                    c.start = RowInfo(**c.start)
                if isinstance(c.end, dict):
                    c.end = RowInfo(**c.end)
            max_line = max(max_line, c.end.row)
            for i in range(c.start.row, c.end.row+1):
                cover_map[i] = 1
        for c in data.not_covered:
            if isinstance(c, dict):
                c = LineCoverage(**c)
                if isinstance(c.start, dict):
                    c.start = RowInfo(**c.start)
                if isinstance(c.end, dict):
                    c.end = RowInfo(**c.end)
            max_line = max(max_line, c.end.row)
            for i in range(c.start.row, c.end.row+1):
                cover_map[i] = 0
        for i in range(1, max_line+1):
            if cover_map.get(i) is None:
                continue
            ET.SubElement(file_data, "lineToCover", attrib={
                "lineNumber": str(i),
                "covered": "true" if (cover_map[i] == 1) else "false",
            })
    return coverage_xml


if __name__ == "__main__":
    # init args
    parser = argparse.ArgumentParser(
        prog='opa-coverage-to-sonarcloud',
        description='Convert opa coverage report to SonarCloud format')
    parser.add_argument("input", help='input json file')
    parser.add_argument("output", help='output xml file')
    args = parser.parse_args()
    with open(args.input, encoding="utf-8") as fp:
        content = json.load(fp)
    overall_coverage = OPACoverage(**content)
    sonarcloud_xml = generate_sonarcloud_xml(overall_coverage)

    tree = ET.ElementTree(sonarcloud_xml)
    ET.indent(tree)
    tree.write(args.output, encoding="utf-8", xml_declaration=True)
