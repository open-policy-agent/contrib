package iptables

import (
	"bytes"
	"testing"
)

func TestAddParam(t *testing.T) {

	var testCase = []struct {
		param  string
		flag   string
		result string
	}{
		{
			"tcp",
			"-p",
			" -p tcp",
		},
		{
			"192.168.0.1",
			"-s",
			" -s 192.168.0.1",
		}, {
			"8080",
			"--sport",
			" --sport 8080",
		}, {
			"192.168.0.1",
			"-d",
			" -d 192.168.0.1",
		},
		{
			"eth0",
			"-i",
			" -i eth0",
		},
		{
			"lo",
			"-o",
			" -o lo",
		},
		{
			"!8080",
			"--sport",
			" ! --sport 8080",
		},
	}

	for _, tt := range testCase {
		var buf bytes.Buffer
		addParam(&buf, tt.param, tt.flag)
		if tt.result != buf.String() {
			t.Errorf("Expected %s , got %s", tt.result, buf.String())
		}
	}
}

func TestAddParams(t *testing.T) {

	var testCase = []struct {
		params []string
		flag   string
		result string
	}{
		{
			[]string{"comment", "contrack"},
			"-m",
			" -m comment -m contrack",
		},
	}

	for _, tt := range testCase {
		var buf bytes.Buffer
		addParams(&buf, tt.params, tt.flag)
		if tt.result != buf.String() {
			t.Errorf("Expected %s , got %s", tt.result, buf.String())
		}
	}
}

func TestAddCommemt(t *testing.T) {
	var testCase = []struct {
		comment string
		result  string
	}{
		{
			"rule for blocking port 8080",
			" -m comment --comment \"rule for blocking port 8080\"",
		},
	}
	for _, tt := range testCase {
		var buf bytes.Buffer
		addComment(&buf, tt.comment)
		if tt.result != buf.String() {
			t.Errorf("Expected %s , got %s", tt.result, buf.String())
		}
	}
}

func TestTCPFlags(t *testing.T) {
	var testCase = []struct {
		tf     tcpFlags
		result string
	}{
		{
			tcpFlags{Flags: []string{"SYN", "ACK", "FIN", "RST"}, FlagsSet: []string{"SYN"}},
			" --tcp-flags SYN,ACK,FIN,RST SYN",
		},
		{
			tcpFlags{Flags: []string{}, FlagsSet: []string{}},
			"",
		},
		{
			tcpFlags{Flags: []string{"SYN", "ACK"}, FlagsSet: []string{"ACK"}},
			" --tcp-flags SYN,ACK ACK",
		},
	}
	for _, tt := range testCase {
		var buf bytes.Buffer
		addTCPFlags(&buf, tt.tf)
		if tt.result != buf.String() {
			t.Errorf("Expected %s , got %s", tt.result, buf.String())
		}
	}
}

func TestAddIPRange(t *testing.T) {
	var testCase = []struct {
		match            []string
		sourceRange      string
		destinationRange string
		result           string
	}{
		{
			[]string{},
			"192.168.1.100-192.168.1.199",
			"",
			" -m iprange --src-range 192.168.1.100-192.168.1.199",
		},
		{
			[]string{},
			"",
			"192.168.1.100-192.168.1.199",
			" -m iprange --dst-range 192.168.1.100-192.168.1.199",
		},
		{
			[]string{"tcp", "iprange"},
			"192.168.1.100-192.168.1.199",
			"192.168.1.100-192.168.1.199",
			" --src-range 192.168.1.100-192.168.1.199 --dst-range 192.168.1.100-192.168.1.199",
		},
		{
			[]string{},
			"192.168.1.100-192.168.1.199",
			"192.168.1.100-192.168.1.199",
			" -m iprange --src-range 192.168.1.100-192.168.1.199 --dst-range 192.168.1.100-192.168.1.199",
		},
	}
	for _, tt := range testCase {
		var buf bytes.Buffer
		addIPRange(&buf, tt.match, tt.sourceRange, tt.destinationRange)
		if tt.result != buf.String() {
			t.Errorf("Expected %s , got %s", tt.result, buf.String())
		}
	}
}

func TestAddCTState(t *testing.T) {
	var testCase = []struct {
		match  []string
		states []string
		result string
	}{
		{
			[]string{"state"},
			[]string{"NEW", "ESTABLISHED", "INVALID"},
			" --ctstate NEW,ESTABLISHED,INVALID",
		},
		{
			[]string{"conntrack", "state"},
			[]string{"NEW", "ESTABLISHED", "INVALID"},
			" --ctstate NEW,ESTABLISHED,INVALID",
		},
		{
			[]string{""},
			[]string{"NEW", "ESTABLISHED", "INVALID"},
			" -m conntrack --ctstate NEW,ESTABLISHED,INVALID",
		},
	}
	for _, tt := range testCase {
		var buf bytes.Buffer
		addCTState(&buf, tt.match, tt.states)
		if tt.result != buf.String() {
			t.Errorf("Expected %s , got %s", tt.result, buf.String())
		}
	}
}

func TestRuleConstruction(t *testing.T) {
	var testcases = []struct {
		rule   Rule
		result string
	}{
		{
			Rule{
				Protocol:        "tcp",
				DestinationPort: "8080",
				Comment:         "block all incoming traffic to port 8080",
				Jump:            "DROP",
			},
			" -p tcp --dport 8080 -j DROP -m comment --comment \"block all incoming traffic to port 8080\"",
		},
	}

	for _, tt := range testcases {
		if tt.result != tt.rule.Construct() {
			t.Errorf("Expected %s,but got %s", tt.result, tt.rule.Construct())
		}
	}
}