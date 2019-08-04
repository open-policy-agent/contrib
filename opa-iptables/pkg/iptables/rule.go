package iptables

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	goiptables "github.com/coreos/go-iptables/iptables"
)

// Rule reperesents an IPTable rule
type Rule struct {
	// This option specifies the packet matching table which the command should operate on.
	// If the kernel is configured with automatic module loading, an attempt will be made to
	// load the appropriate module for that table if it is not already there.
	// Choices - filter | nat | mangle | raw | security
	// Default - filter
	Table string `json:"table,omitempty"`

	// Specify the iptables chain to modify.
	// This could be a user-defined chain or one of the standard iptables chains,
	// like INPUT, FORWARD, OUTPUT, PREROUTING, POSTROUTING, SECMARK or CONNSECMARK.
	Chain string `json:"chain,omitempty"`

	// Destination port or port range specification. This can either be a service name or a port number.
	// An inclusive range can also be specified, using the format first:last.
	// If the first port is omitted, '0' is assumed; if the last is omitted, '65535' is assumed.
	// If the first port is greater than the second one they will be swapped.
	// This is only valid if the rule also specifies one of the following protocols: tcp, udp, dccp or sctp.
	DestinationPort string `json:"destination_port,omitempty"`

	// Destination Address specification.
	// Address can be either a network name, a hostname, a network IP address (with /mask), or a plain IP address.
	// Hostnames will be resolved once only, before the rule is submitted to the kernel. Please note that specifying
	// any name to be resolved with a remote query such as DNS is a really bad idea.
	// The mask can be either a network mask or a plain number, specifying the number of 1's at the left side of the network mask.
	// Thus, a mask of 24 is equivalent to 255.255.255.0. A ! argument before the address specification inverts the sense of the address.
	DestinationAddress string `json:"destination,omitempty"`

	// Specifies the destination IP range to match in the iprange module.
	DestinationRange string `json:"dst_range,omitempty"`

	// Source Address specification.
	// Address can be either a network name, a hostname, a network IP address (with /mask), or a plain IP address.
	// Hostnames will be resolved once only, before the rule is submitted to the kernel.
	// Please note that specifying any name to be resolved with a remote query such as DNS is a really bad idea.
	// The mask can be either a network mask or a plain number, specifying the number of 1's at the left side of the network mask.
	// Thus, a mask of 24 is equivalent to 255.255.255.0. A ! argument before the address specification inverts the sense of the address.
	SourceAddress string `json:"source,omitempty"`

	// Source port or port range specification.
	// This can either be a service name or a port number.
	// An inclusive range can also be specified, using the format first:last.
	// If the first port is omitted, 0 is assumed; if the last is omitted, 65535 is assumed.
	// If the first port is greater than the second one they will be swapped.
	SourcePort string `json:"source_port,omitempty"`

	// Specifies the source IP range to match in the iprange module.
	SourceRange string `json:"src_range,omitempty"`

	// This specifies a destination address to use with DNAT.
	// Without this, the destination address is never altered.
	ToDestination string `json:"to_destination,omitempty"`

	// This specifies a source address to use with SNAT.
	// Without this, the source address is never altered.
	ToSource string `json:"to_source,omitempty"`

	// This specifies a destination port or range of ports to use, without this, the destination port is never altered.
	// This is only valid if the rule also specifies one of the protocol tcp, udp, dccp or sctp.
	ToPorts string `json:"to_ports,omitempty"`

	// Whether the rule should be appended at the bottom or inserted at the top.
	// If the rule already exists the chain will not be modified.
	// Choices : append | insert
	// Defualt : append
	Action string `json:"action,omitempty"`

	//Insert the rule as the given rule number.
	// This works only with action=insert.
	RuleNumber string `json:"rule_num,omitempty"`

	// This specifies the target of the rule; i.e., what to do if the packet matches it.
	// The target can be a user-defined chain (other than the one this rule is in), one of
	// the special builtin targets which decide the fate of the packet immediately, or an extension (see EXTENSIONS for more at http://ipset.netfilter.org/iptables-extensions.man.html).
	// If this option is omitted in a rule (and -g is not used), then matching the rule will have no effect
	// on the packet's fate, but the counters on the rule will be incremented.
	Jump string `json:"jump,omitempty"`

	// Name of an interface via which a packet was received (only for packets entering the INPUT, FORWARD and PREROUTING chains).
	// When the ! argument is used before the interface name, the sense is inverted.
	// If the interface name ends in a +, then any interface which begins with this name will match.
	// If this option is omitted, any interface name will match.
	InInterface string `json:"in_interface,omitempty"`

	// Name of an interface via which a packet is going to be sent (for packets entering the FORWARD, OUTPUT and POSTROUTING chains).
	// When the ! argument is used before the interface name, the sense is inverted.
	// If the interface name ends in a +, then any interface which begins with this name will match.
	// If this option is omitted, any interface name will match.
	OutInterface string `json:"out_interface,omitempty"`

	// The protocol of the rule or of the packet to check.
	// The specified protocol can be one of tcp, udp, udplite, icmp, esp, ah, sctp or the special keyword all, or it can be a numeric value, representing one of these protocols or a different one.
	// A protocol name from /etc/protocols is also allowed.
	// A ! argument before the protocol inverts the test.
	// The number zero is equivalent to all.
	// all will match with all protocols and is taken as default when this option is omitted.
	Protocol string `json:"protocol,omitempty"`

	// TCP flags specification.
	// tcp_flags expects a struct with the two keys flags and flags_set.
	TCPFlags TcpFlags `json:"tcp_flags,omitempty"`

	// ctstate is a list of the connection states to match in the conntrack module.
	// Possible states are INVALID, NEW, ESTABLISHED, RELATED, UNTRACKED, SNAT, DNAT
	Ctstate []string `json:"ctstate,omitempty"`

	// Specifies a match to use, that is, an extension module that tests for a specific property.
	// The set of matches make up the condition under which a target is invoked.
	// Matches are evaluated first to last if specified as an array and work in short-circuit fashion,
	// i.e. if one extension yields false, evaluation will stop.
	Match []string `json:"match,omitempty"`

	LogPrefix string `json:"log_prefix,omitempty"`

	// This specifies a comment that will be added to the rule.
	Comment string `json:"comment,omitempty"`
}

type TcpFlags struct {
	// List of flags you want to examine.
	Flags []string `json:"flags,omitempty"`
	// Flags to be set.
	FlagsSet []string `json:"flags_set,omitempty"`
}

type ruleSpec struct {
	spec []string
}

type RuleSet struct {
	Rules []Rule `json:"result"`
}

func (rs *ruleSpec) addParam(param string, flag string) {
	if param != "" {
		if param[0] == '!' {
			rs.spec = append(rs.spec, "!", flag, param[1:])
		} else {
			rs.spec = append(rs.spec, flag, param)
		}
	}
}

func (rs *ruleSpec) addParams(params []string, flag string) {
	if len(params) > 0 {
		rs.spec = append(rs.spec, flag)
		for _, param := range params {
			rs.spec = append(rs.spec, param)
		}
	}
}

func (rs *ruleSpec) addMatch(param string) {
	if param != "" {
		rs.spec = append(rs.spec, "-m", param)
	}
}

func (rs *ruleSpec) addComment(matchs []string, comment string) {
	for _, match := range matchs {
		if match == "comment" && comment != "" {
			rs.addMatch("comment")
			rs.spec = append(rs.spec, "--comment", fmt.Sprintf("\"%s\"", comment))
			return
		}
	}
	if comment != "" {
		rs.addMatch("comment")
		rs.spec = append(rs.spec, "--comment", fmt.Sprintf("\"%s\"", comment))
	}
}

func (rs *ruleSpec) addIPRange(matchs []string, sourceRange, destinationRange string) {
	for _, match := range matchs {
		if match == "iprange" {
			rs.addMatch("iprange")
			rs.addParam(sourceRange, "--src-range")
			rs.addParam(destinationRange, "--dst-range")
			return
		}
	}
	if sourceRange != "" || destinationRange != "" {
		rs.addMatch("iprange")
		rs.addParam(sourceRange, "--src-range")
		rs.addParam(destinationRange, "--dst-range")
	}
}

func (rs *ruleSpec) addTCPFlags(tf TcpFlags) {
	if len(tf.Flags) > 0 && len(tf.FlagsSet) > 0 {
		rs.addParams([]string{strings.Join(tf.Flags, ","), strings.Join(tf.FlagsSet, ",")}, "--tcp-flags")
	}
}

func (rs *ruleSpec) addCTState(matchs, states []string) {
	for _, match := range matchs {
		if match == "conntrack" {
			rs.addMatch("conntrack")
			rs.addParam(strings.Join(states, ","), "--ctstate")
			return
		} else if match == "state" {
			rs.addMatch("conntrack")
			rs.addParam(strings.Join(states, ","), "--ctstate")
			return
		}
	}
	if len(states) > 0 && !isEmpty(states) {
		rs.addMatch("conntrack")
		rs.addParam(strings.Join(states, ","), "--ctstate")
	}
}

func isEmpty(arr []string) bool {
	for _, v := range arr {
		if v != "" {
			return false
		}
	}
	return true
}

// Construct IPTable rule from struct
func (r *Rule) Construct() []string {
	var rs ruleSpec
	rs.addParam(r.Protocol, "-p")
	rs.addParam(r.SourceAddress, "-s")
	rs.addParam(r.SourcePort, "--sport")
	rs.addParam(r.DestinationAddress, "-d")
	rs.addParam(r.DestinationPort, "--dport")
	rs.addParam(r.InInterface, "-i")
	rs.addParam(r.OutInterface, "-o")
	rs.addIPRange(r.Match, r.SourceRange, r.DestinationRange)
	rs.addCTState(r.Match, r.Ctstate)
	rs.addTCPFlags(r.TCPFlags)
	rs.addParam(r.Jump, "-j")
	rs.addParam(r.ToSource, "--to-source")
	rs.addParam(r.ToDestination, "--to-destination")
	rs.addParam(r.ToPorts, "--to-ports")
	rs.addParam(r.LogPrefix, "--log-prefix")
	rs.addComment(r.Match, r.Comment)
	return rs.spec
}

func (r Rule) String() string {
	spec := []string{r.Table, r.Chain}
	spec = append(spec, r.Construct()...)
	return strings.Join(spec, " ")
}

func UnmarshalRules(rules []byte) (RuleSet, error) {
	var rs RuleSet
	err := json.Unmarshal(rules, &rs)
	if err != nil {
		return RuleSet{}, err
	}
	for i := range rs.Rules {
		rs.Rules[i].init()
	}
	return rs, nil
}

func (r *Rule) AddRule() error {
	ipt, err := goiptables.NewWithProtocol(goiptables.ProtocolIPv4)
	if err != nil {
		return err
	}

	switch r.Action {
	// inserts rulespec to specified table/chain (in specified position)
	case "insert":
		if r.RuleNumber != "" {
			ruleNum, err := strconv.Atoi(r.RuleNumber)
			if err != nil {
				return err
			}
			ipt.Insert(r.Table, r.Chain, ruleNum, r.Construct()...)
		} else {
			return errors.New("to use insert action ,you must need to provides rule_number")
		}
	default:
		// appends rulespec to specified table/chain
		return ipt.AppendUnique(r.Table, r.Chain, r.Construct()...)
	}
	return nil
}

func (r *Rule) DeleteRule() error {
	ipt, err := goiptables.NewWithProtocol(goiptables.ProtocolIPv4)
	if err != nil {
		return err
	}
	return ipt.Delete(r.Table, r.Chain, r.Construct()...)
}

// adding default values to IPTables rules (if user not provides it)
func (r *Rule) init() {
	r.Table = strings.ToLower(r.Table)
	if r.Table == "" {
		r.Table = "filter"
	}
	r.Chain = strings.ToUpper(r.Chain)
	if r.Chain == "" {
		r.Chain = "INPUT"
	}
}

func ListRules(table, chain string) ([]string, error) {
	ipt, err := goiptables.NewWithProtocol(goiptables.ProtocolIPv4)
	if err != nil {
		return nil, err
	}
	return ipt.List(strings.ToLower(table), strings.ToUpper(chain))
}
