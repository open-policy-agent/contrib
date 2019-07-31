package converter

import (
	"strings"
	"testing"
)

var rule = `iptables -t filter -A INPUT -m conntrack,tcp,comment -p tcp -i eth0 -o eth0 -s 192.168.0.1 -d 127.0.0.1 --ctstate ESTABLISHED,RELATED -dport 8080 -sport 9090 -j DROP --to-ports 80 --tcp-flags ALL ACK,RST,SYN,FIN --comment "hello world"`

var expected = []string{`{
    "table": "filter",
    "chain": "INPUT",
    "destination_port": "8080",
    "destination": "127.0.0.1",
    "source": "192.168.0.1",
    "source_port": "9090",
    "to_ports": "80",
    "jump": "DROP",
    "in_interface": "eth0",
    "out_interface": "eth0",
    "protocol": "tcp",
    "tcp_flags": {
        "flags": [
            "ALL"
        ],
        "flags_set": [
            "ACK",
            "RST",
            "SYN",
            "FIN"
        ]
    },
    "ctstate": [
        "ESTABLISHED",
        "RELATED"
    ],
    "match": [
        "conntrack",
        "tcp",
        "comment"
    ],
    "comment": "hello world"
}`}

func TestIPTableToJSON(t *testing.T) {
	rules, err := IPTableToJSON(strings.NewReader(rule))
	if err != nil {
		t.Error(err)
	}
	for i, rule := range rules {
        if rule != expected[i] {
            t.Errorf("wanted: %v, but got: %v",expected[i],rule)
        } 
    }
}