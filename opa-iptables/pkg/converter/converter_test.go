package converter

import (
	"strings"
	"testing"
)

var file = 
`iptables -t filter -A INPUT -m conntrack,tcp,comment -p tcp -i eth0 -o eth0 -s 192.168.0.1 -d 127.0.0.1 -dport 8080 -sport 9090 -j DROP --to-ports 80 --tcp-flags ALL ACK,RST,SYN,FIN --comment "hello world"
iptables -A INPUT -i lo -j ACCEPT
iptables -A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
iptables -A INPUT -s 15.15.15.51 -j DROP
iptables -A INPUT -p tcp --dport 22 -m conntrack --ctstate NEW,ESTABLISHED -j ACCEPT
iptables -A INPUT -p tcp -s 15.15.15.0/24 --dport 22 -m conntrack --ctstate NEW,ESTABLISHED -j ACCEPT
iptables -A INPUT -p tcp -s 15.15.15.0/24 --dport 3306 -m conntrack --ctstate NEW,ESTABLISHED -j ACCEPT
iptables -A OUTPUT -o eth1 -p tcp --sport 5432 -m conntrack --ctstate ESTABLISHED -j ACCEPT
iptables -A INPUT -p tcp --dport 143 -m conntrack --ctstate NEW,ESTABLISHED -j ACCEPT`

var expected = 
`[{
    "table": "filter",
    "chain": "INPUT",
    "destination_port": "127.0.0.1",
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
        ""
    ],
    "match": [
        "conntrack",
        "tcp",
        "comment"
    ],
    "comment": "hello world"
},
{
    "chain": "INPUT",
    "jump": "ACCEPT",
    "in_interface": "lo",
    "tcp_flags": {},
    "ctstate": [
        ""
    ],
    "match": [
        ""
    ]
},
{
    "chain": "INPUT",
    "jump": "ACCEPT",
    "tcp_flags": {},
    "ctstate": [
        "ESTABLISHED",
        "RELATED"
    ],
    "match": [
        "conntrack"
    ]
},
{
    "chain": "INPUT",
    "source": "15.15.15.51",
    "jump": "DROP",
    "tcp_flags": {},
    "ctstate": [
        ""
    ],
    "match": [
        ""
    ]
},
{
    "chain": "INPUT",
    "jump": "ACCEPT",
    "protocol": "tcp",
    "tcp_flags": {},
    "ctstate": [
        "NEW",
        "ESTABLISHED"
    ],
    "match": [
        "conntrack"
    ]
},
{
    "chain": "INPUT",
    "source": "15.15.15.0/24",
    "jump": "ACCEPT",
    "protocol": "tcp",
    "tcp_flags": {},
    "ctstate": [
        "NEW",
        "ESTABLISHED"
    ],
    "match": [
        "conntrack"
    ]
},
{
    "chain": "INPUT",
    "source": "15.15.15.0/24",
    "jump": "ACCEPT",
    "protocol": "tcp",
    "tcp_flags": {},
    "ctstate": [
        "NEW",
        "ESTABLISHED"
    ],
    "match": [
        "conntrack"
    ]
},
{
    "chain": "OUTPUT",
    "source_port": "5432",
    "jump": "ACCEPT",
    "out_interface": "eth1",
    "protocol": "tcp",
    "tcp_flags": {},
    "ctstate": [
        "ESTABLISHED"
    ],
    "match": [
        "conntrack"
    ]
},
{
    "chain": "INPUT",
    "jump": "ACCEPT",
    "protocol": "tcp",
    "tcp_flags": {},
    "ctstate": [
        "NEW",
        "ESTABLISHED"
    ],
    "match": [
        "conntrack"
    ]
}]`

func TestIPTableToJSON(t *testing.T) {
	rules, err := IPTableToJSON(strings.NewReader(file))
	if err != nil {
		t.Error(err)
	}
	if rules != expected {
		t.Errorf("wanted: %v, got: %v", expected, rules)
	}
}
