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


def generate_sources() -> ET.Element:
    """<sources>
        <source>/Users/leobalter/dev/testing/solutions/3</source>
    </sources>"""
    sources = ET.Element("sources")
    current_source = ET.Element("source")
    current_source.text = os.getcwd()
    sources.append(current_source)
    return sources


def generate_package(coverage: OPACoverage) -> ET.Element:
    packages = ET.Element("packages")
    package = ET.Element("package", attrib={
        "complexity": "0",
        "line-rate": str(coverage.coverage / 100),
        "name": "",
    })
    classes = ET.Element("classes")
    packages.append(package)
    package.append(classes)
    for path, data in coverage.files.items():
        if isinstance(data, dict):
            data = FileCoverage(**data)
        class_data = ET.Element("class", attrib={
            "complexity": "0",
            "line-rate": str(data.coverage / 100),
            "filename": path,
            "name": path,
        })
        classes.append(class_data)
        ET.SubElement(class_data, "methods")
        line_data = ET.Element("lines")
        class_data.append(line_data)
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
            ET.SubElement(line_data, "line", attrib={
                "number": str(i),
                "hits": str(cover_map[i]),
            })
    return packages


if __name__ == "__main__":
    # init args
    parser = argparse.ArgumentParser(
        prog='opa-coverage-to-cobertura',
        description='Convert opa coverage report to cobertura')
    parser.add_argument("input", help='input json file')
    parser.add_argument("output", help='output xml file')
    args = parser.parse_args()
    with open(args.input, encoding="utf-8") as fp:
        content = json.load(fp)
    overall_coverage = OPACoverage(**content)
    cobertura_xml = ET.Element("coverage", attrib={
        "lines-covered": str(overall_coverage.covered_lines),
        "line-rate": str(overall_coverage.coverage / 100),
        "lines-valid": str(overall_coverage.lines_valid),
        "complexity": "0",
        "version": "0.1",
        "timestamp": str(math.floor(time.time() * 1000)),
    })
    cobertura_xml.append(generate_package(overall_coverage))

    tree = ET.ElementTree(cobertura_xml)
    ET.indent(tree)
    tree.write(args.output, encoding="utf-8", xml_declaration=True)
