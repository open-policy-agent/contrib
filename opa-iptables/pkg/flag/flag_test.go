package flag

import (
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	var testcases = []struct {
		arguments string
		flagSet   IPTableflagSet
		err       error
	}{
		{
			"iptables -t filter -A INPUT -m conntrack,tcp -p tcp -i eth0 -o eth0 -s 192.168.0.1 -d 127.0.0.1 -dport 8080 -sport 9090 -j DROP --to-ports 80 --tcp-flags ALL ACK,RST,SYN,FIN",
			IPTableflagSet{
				TableFlag:    "filter",
				ChainFlag:    "INPUT",
				ProtocolFlag: "tcp",
				DportFlag:    "8080",
				SportFlag:    "9090",
				JumpFlag:     "DROP",
				TCPFlag: TCPFlags{
					Flags:    []string{"ALL"},
					FlagsSet: []string{"ACK", "RST", "SYN", "FIN"},
				},
				InInterfaceFlag:  "eth0",
				OutInterfaceFlag: "eth0",
				SourceFlag:       "192.168.0.1",
				DestinationFlag:  "127.0.0.1",
				ToPortFlag:       "80",
				MatchFlag:        "conntrack,tcp",
			},
			nil,
		},
		{
			"iptables -t -A PREROUTING",
			IPTableflagSet{},
			argumentError{flagName:"t",numArg:1},
		},
		{
			"iptables -notdefined notsure",
			IPTableflagSet{},
			flagError{name:"notdefined"},
		},
	}

	for _, tt := range testcases {
		fs := NewFlagSet("iptables", ContinueOnError)
		var iptFlagset IPTableflagSet
		fs.InitFlagSet(&iptFlagset)
		err := fs.Parse(strings.Split(tt.arguments, " "))
		if err == tt.err {
			t.Errorf("wanted: %v, got: %v",tt.err,err)
		}
		if !reflect.DeepEqual(iptFlagset, tt.flagSet) {
			t.Errorf("wanted: %#v, but got: %#v", tt.flagSet, iptFlagset)
		}
	}
}
