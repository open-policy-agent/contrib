package converter

import (
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
		Table:              strings.ToLower(tf.TableFlag),
		Chain:              strings.ToUpper(tf.ChainFlag),
		Protocol:           tf.ProtocolFlag,
		DestinationPort:    tf.DportFlag,
		DestinationAddress: tf.DestinationFlag,
		SourceAddress:      tf.SourceFlag,
		SourcePort:         tf.SportFlag,
		InInterface:        tf.InInterfaceFlag,
		OutInterface:       tf.OutInterfaceFlag,
		DestinationRange:   tf.DesRangeFlag,
		SourceRange:        tf.SrcRangeFlag,
		Jump:               tf.JumpFlag,
		ToPorts:            tf.ToPortFlag,
		Match:              strings.Split(tf.MatchFlag, ","),
		Ctstate:            strings.Split(tf.CTStateFlag, ","),
		TCPFlags:           iptables.TcpFlags(tf.TCPFlag),
		Comment:            tf.Comment,
	}

	return json.MarshalIndent(r, "", "    ")
}

// IPTableToJSON reads '\n' delimeted rules from reader, parse each rule and returns rules describes in JSON format as a string.
func IPTableToJSON(reader io.Reader) ([]string, error) {
	var jsonRules []string
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	rules := strings.Split(string(b), "\n")
	for _, rule := range rules {
		fs := flag.NewFlagSet("iptables", flag.ContinueOnError)
		var tableFlagset flag.IPTableflagSet
		fs.InitFlagSet(&tableFlagset)

		// Parse line as a shell words
		// i.e "iptables --comment "hello world""
		// args should be ["iptables","--comment","hello world"]
		args, err := shellwords.Parse(rule)
		if err != nil {
			jsonRules = append(jsonRules, "\"Error: "+err.Error()+"\"")
			continue
		}

		err = fs.Parse(args)
		if err != nil {
			jsonRules = append(jsonRules, "\"Error: "+err.Error()+"\"")
			continue
		}

		rule, err := marshal(tableFlagset)
		if err != nil {
			jsonRules = append(jsonRules, "\"Error: "+err.Error()+"\"")
			continue
		}

		jsonRules = append(jsonRules, string(rule))
	}
	return jsonRules, nil
}
