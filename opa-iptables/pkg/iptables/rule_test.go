package iptables

import (
	"reflect"
	"testing"
)

func TestAddParam(t *testing.T) {

	var testCase = []struct {
		param  string
		flag   string
		result []string
	}{
		{
			"tcp",
			"-p",
			[]string{"-p" ,"tcp"},
		},
		{
			"192.168.0.1",
			"-s",
			[]string{"-s" ,"192.168.0.1"},
		}, {
			"8080",
			"--sport",
			[]string{"--sport","8080" },
		}, {
			"192.168.0.1",
			"-d",
			[]string{"-d", "192.168.0.1"},
		},
		{
			"eth0",
			"-i",
			[]string{"-i" ,"eth0"},
		},
		{
			"!8080",
			"--sport",
			[]string{"!", "--sport", "8080"},
		},
	}

	for _, tt := range testCase {
		var rs ruleSpec
		rs.addParam(tt.param, tt.flag)
		if !reflect.DeepEqual(rs.spec,tt.result) {
			t.Errorf("Expected %s , got %s", tt.result, rs.spec)
		}
	}
}

func TestAddParams(t *testing.T) {

	var testCase = []struct {
		params []string
		flag   string
		result []string
	}{
		{
			[]string{"comment", "contrack"},
			"-m",
			[]string{"-m", "comment","contrack"},
		},
	}

	for _, tt := range testCase {
		var rs ruleSpec
		rs.addParams(tt.params, tt.flag)
		if !reflect.DeepEqual(rs.spec,tt.result) {
			t.Errorf("Expected %s , got %s", tt.result, rs.spec)
		}
	}
}

func TestAddCommemt(t *testing.T) {
	var testCase = []struct {
		comment string
		result  []string
	}{
		{
			"rule for blocking port 8080",
			[]string{"-m", "comment" ,"--comment", "\"rule for blocking port 8080\""},
		},
	}
	for _, tt := range testCase {
		var rs ruleSpec
		rs.addComment(tt.comment)
		if !reflect.DeepEqual(rs.spec,tt.result) {
			t.Errorf("Expected %s , got %s", tt.result, rs.spec)
		}
	}
}

func TestTCPFlags(t *testing.T) {
	var testCase = []struct {
		tf     tcpFlags
		result []string
	}{
		{
			tcpFlags{Flags: []string{"SYN", "ACK", "FIN", "RST"}, FlagsSet: []string{"SYN"}},
			[]string{"--tcp-flags", "SYN,ACK,FIN,RST","SYN"},
		},
		{
			tcpFlags{Flags: []string{"SYN", "ACK"}, FlagsSet: []string{"ACK"}},
			[]string{"--tcp-flags" ,"SYN,ACK", "ACK"},
		},
	}
	for _, tt := range testCase {
		var rs ruleSpec
		rs.addTCPFlags(tt.tf)
		if !reflect.DeepEqual(tt.result,rs.spec) {
			t.Errorf("Expected %s , got %s", tt.result, rs.spec)
		}
	}
}

func TestAddIPRange(t *testing.T) {
	var testCase = []struct {
		match            []string
		sourceRange      string
		destinationRange string
		result           []string
	}{
		{
			[]string{},
			"192.168.1.100-192.168.1.199",
			"",
			[]string{"-m", "iprange" ,"--src-range", "192.168.1.100-192.168.1.199"},
		},
		{
			[]string{},
			"",
			"192.168.1.100-192.168.1.199",
			[]string{"-m" ,"iprange" ,"--dst-range", "192.168.1.100-192.168.1.199"},
		},
		{
			[]string{"tcp", "iprange"},
			"192.168.1.100-192.168.1.199",
			"192.168.1.100-192.168.1.199",
			[]string{"--src-range", "192.168.1.100-192.168.1.199", "--dst-range", "192.168.1.100-192.168.1.199"},
		},
		{
			[]string{},
			"192.168.1.100-192.168.1.199",
			"192.168.1.100-192.168.1.199",
			[]string{"-m", "iprange", "--src-range", "192.168.1.100-192.168.1.199", "--dst-range", "192.168.1.100-192.168.1.199"},
		},
	}
	for _, tt := range testCase {
		var rs ruleSpec
		rs.addIPRange(tt.match, tt.sourceRange, tt.destinationRange)
		if !reflect.DeepEqual(rs.spec,tt.result) {
			t.Errorf("Expected %s , got %s", tt.result, rs.spec)
		}
	}
}

func TestAddCTState(t *testing.T) {
	var testCase = []struct {
		match  []string
		states []string
		result []string
	}{
		{
			[]string{"state"},
			[]string{"NEW", "ESTABLISHED", "INVALID"},
			[]string{"--ctstate", "NEW,ESTABLISHED,INVALID"},
		},
		{
			[]string{"conntrack", "state"},
			[]string{"NEW", "ESTABLISHED", "INVALID"},
			[]string{"--ctstate", "NEW,ESTABLISHED,INVALID"},
		},
		{
			[]string{""},
			[]string{"NEW", "ESTABLISHED", "INVALID"},
			[]string{"-m", "conntrack" ,"--ctstate", "NEW,ESTABLISHED,INVALID"},
		},
	}
	for _, tt := range testCase {
		var rs ruleSpec
		rs.addCTState(tt.match, tt.states)
		if !reflect.DeepEqual(rs.spec,tt.result) {
			t.Errorf("Expected %s , got %s", tt.result, rs.spec)
		}
	}
}

func TestRuleConstruction(t *testing.T) {
	var testcases = []struct {
		rule   Rule
		result []string
	}{
		{
			Rule{
				Table:           "filter",
				Chain:           "INPUT",
				Protocol:        "tcp",
				DestinationPort: "8080",
				Comment:         "block all incoming traffic to port 8080",
				Jump:            "DROP",
			},
			[]string{"-p", "tcp", "--dport", "8080", "-j", "DROP" ,"-m", "comment" , "--comment" , "\"block all incoming traffic to port 8080\""},
		},
		{
			Rule{
				Table:           "nat",
				Chain:           "PREROUTING",
				Protocol:        "tcp",
				InInterface:     "eth0",
				DestinationPort: "80",
				ToPorts:         "8080",
				Jump:            "REDIRECT",
				Comment:         "Redirect web traffic from port 80 to port 8080",
			},
			[]string{"-p", "tcp", "--dport", "80", "-i", "eth0", "-j", "REDIRECT","--to-ports", "8080","-m", "comment", "--comment", "\"Redirect web traffic from port 80 to port 8080\""},
		},
		{
			Rule{
				Table: "filter",
				Chain: "OUTPUT",
				Protocol: "tcp",
				TCPFlags: tcpFlags{
					Flags:[]string{"ACK","RST","SYN","FIN"},
					FlagsSet:[]string{"SYN"},
				},
				Jump: "DROP",
			},
			[]string{"-p","tcp","--tcp-flags","ACK,RST,SYN,FIN","SYN","-j","DROP"},
		},
	}

	for _, tt := range testcases {
		if !reflect.DeepEqual(tt.result,tt.rule.Construct()) {
			t.Errorf("Expected %s,but got %s", tt.result, tt.rule.Construct())
		}
	}
}


func TestAddRule(t *testing.T) {
	var testcases = []struct{
		rule Rule
	}{
		{
			Rule{
				Table:           "filter",
				Chain:           "INPUT",
				Protocol:        "tcp",
				DestinationPort: "8080",
				Comment:         "block all incoming traffic to port 8080",
				Jump:            "DROP",
			},
		},
		{
			Rule{
				Table:           "nat",
				Chain:           "PREROUTING",
				Protocol:        "tcp",
				InInterface:     "eth0",
				DestinationPort: "80",
				ToPorts:         "8080",
				Jump:            "REDIRECT",
				Comment:         "Redirect web traffic from port 80 to port 8080",
			},
		},
		{
			Rule{
				Table: "filter",
				Chain: "OUTPUT",
				Protocol: "tcp",
				TCPFlags: tcpFlags{
					Flags:[]string{"ALL"},
					FlagsSet:[]string{"ACK","RST","SYN","FIN"},
				},
				Jump: "DROP",
			},
		},
	}

	for _, tt := range testcases {
		err := tt.rule.AddRule() 
		if err != nil {
			t.Errorf("Got error %s",err)
		}
	}
}