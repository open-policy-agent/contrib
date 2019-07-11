package converter

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"

	"github.com/mattn/go-shellwords"
	"github.com/open-policy-agent/contrib/opa-iptables/pkg/flag"
	"github.com/open-policy-agent/contrib/opa-iptables/pkg/iptables"
)

func marshal(tf flag.IPTableflagSet) ([]byte, error) {
	r := iptables.Rule{
		Table:            tf.TableFlag,
		Chain:            tf.ChainFlag,
		Protocol:         tf.ProtocolFlag,
		SourceAddress:    tf.SourceFlag,
		DestinationPort:  tf.DestinationFlag,
		SourcePort:       tf.SportFlag,
		InInterface:      tf.InInterfaceFlag,
		OutInterface:     tf.OutInterfaceFlag,
		DestinationRange: tf.DesRangeFlag,
		SourceRange:      tf.SrcRangeFlag,
		Jump:             tf.JumpFlag,
		ToPorts:          tf.ToPortFlag,
		Match:            strings.Split(tf.MatchFlag, ","),
		Ctstate:          strings.Split(tf.CTStateFlag, ","),
		TCPFlags:         iptables.TcpFlags(tf.TCPFlag),
		Comment:          tf.Comment,
	}

	return json.MarshalIndent(r, "", "    ")
}

// IPTableToJSON reads '\n' delimeted rules from reader, parse each rule and returns rules describes in JSON format as a string.
func IPTableToJSON(reader io.Reader) (string, error) {
	var buf bytes.Buffer
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	rules := strings.Split(string(b), "\n")
	for i, rule := range rules {
		fs := flag.NewFlagSet("iptables", flag.ContinueOnError)
		var tableFlagset flag.IPTableflagSet
		fs.InitFlagSet(&tableFlagset)

		// Parse line as a shell words
		// i.e "iptables --comment "hello world""
		// args should be ["iptables","--comment","hello world"]
		args, err := shellwords.Parse(rule)
		if err != nil {
			return "", err
		}

		err = fs.Parse(args)
		if err != nil {
			return "", err
		}

		rule, err := marshal(tableFlagset)
		if err != nil {
			return "", err
		}

		// for first rule
		if i == 0 {
			buf.Write([]byte("["))
		}

		// add each rule to buffer
		buf.Write(rule)

		// add ',' and '\n' to each rule except last rule
		if i != len(rules)-1 {
			buf.WriteByte(',')
			buf.WriteByte('\n')
		}
		// for last rule
		if i == len(rules)-1 {
			buf.Write([]byte("]"))
		}
	}
	return buf.String(), nil
}
