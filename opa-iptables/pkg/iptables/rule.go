package iptables

// Rule reperesents an IPTable rule
type Rule struct {
	// This option specifies the packet matching table which the command should operate on.
	// If the kernel is configured with automatic module loading, an attempt will be made to
	// load the appropriate module for that table if it is not already there.
	// Choices - filter | nat | mangle | raw | security
	// Default - filter
	Table string `json:"table"`

	// Specify the iptables chain to modify.
	// This could be a user-defined chain or one of the standard iptables chains,
	// like INPUT, FORWARD, OUTPUT, PREROUTING, POSTROUTING, SECMARK or CONNSECMARK.
	Chain string `json:"chain"`

	// This specifies a comment that will be added to the rule.
	Comment string `json:"comment"`

	// Destination port or port range specification. This can either be a service name or a port number.
	// An inclusive range can also be specified, using the format first:last.
	// If the first port is omitted, '0' is assumed; if the last is omitted, '65535' is assumed.
	// If the first port is greater than the second one they will be swapped.
	// This is only valid if the rule also specifies one of the following protocols: tcp, udp, dccp or sctp.
	DestinationPort string `json:"destination_port"`

	// Destination Address specification.
	// Address can be either a network name, a hostname, a network IP address (with /mask), or a plain IP address.
	// Hostnames will be resolved once only, before the rule is submitted to the kernel. Please note that specifying
	// any name to be resolved with a remote query such as DNS is a really bad idea.
	// The mask can be either a network mask or a plain number, specifying the number of 1's at the left side of the network mask.
	// Thus, a mask of 24 is equivalent to 255.255.255.0. A ! argument before the address specification inverts the sense of the address.
	DestinationAddress string `json:"destination"`

	// Specifies the destination IP range to match in the iprange module.
	DestinationRange string `json:"dst_range"`

	// Source Address specification.
	// Address can be either a network name, a hostname, a network IP address (with /mask), or a plain IP address.
	// Hostnames will be resolved once only, before the rule is submitted to the kernel.
	// Please note that specifying any name to be resolved with a remote query such as DNS is a really bad idea.
	// The mask can be either a network mask or a plain number, specifying the number of 1's at the left side of the network mask.
	// Thus, a mask of 24 is equivalent to 255.255.255.0. A ! argument before the address specification inverts the sense of the address.
	SourceAddress string `json:"source"`

	// Source port or port range specification.
	// This can either be a service name or a port number.
	// An inclusive range can also be specified, using the format first:last.
	// If the first port is omitted, 0 is assumed; if the last is omitted, 65535 is assumed.
	// If the first port is greater than the second one they will be swapped.
	SourcePort string `json:"source_port"`

	// Specifies the source IP range to match in the iprange module.
	SourceRange string `json:"src_range"`

	// This specifies a destination address to use with DNAT.
	// ithout this, the destination address is never altered.
	ToDestination string `json:"to_destination"`

	// This specifies a source address to use with SNAT.
	// Without this, the source address is never altered.
	ToSource string `json:"to_source"`

	// This specifies a destination port or range of ports to use, without this, the destination port is never altered.
	// This is only valid if the rule also specifies one of the protocol tcp, udp, dccp or sctp.
	ToPorts string `json:"to_ports"`

	// Whether the rule should be appended at the bottom or inserted at the top.
	// If the rule already exists the chain will not be modified.
	// Choices : append | insert
	// Defualt : append
	Action string `json:"action"`

	//Insert the rule as the given rule number.
	// This works only with action=insert.
	RuleNumber string `json:"rule_num"`

	// This specifies the target of the rule; i.e., what to do if the packet matches it.
	// The target can be a user-defined chain (other than the one this rule is in), one of
	// the special builtin targets which decide the fate of the packet immediately, or an extension (see EXTENSIONS for more at http://ipset.netfilter.org/iptables-extensions.man.html).
	// If this option is omitted in a rule (and -g is not used), then matching the rule will have no effect
	// on the packet's fate, but the counters on the rule will be incremented.
	Jump string `json:"jump"`

	// Name of an interface via which a packet was received (only for packets entering the INPUT, FORWARD and PREROUTING chains).
	// When the ! argument is used before the interface name, the sense is inverted.
	// If the interface name ends in a +, then any interface which begins with this name will match.
	// If this option is omitted, any interface name will match.
	InInterface string `json:"in_interface"`

	// Name of an interface via which a packet is going to be sent (for packets entering the FORWARD, OUTPUT and POSTROUTING chains).
	// When the ! argument is used before the interface name, the sense is inverted.
	// If the interface name ends in a +, then any interface which begins with this name will match.
	// If this option is omitted, any interface name will match.
	OutInterface string `json:"out_interface"`

	// The protocol of the rule or of the packet to check.
	// The specified protocol can be one of tcp, udp, udplite, icmp, esp, ah, sctp or the special keyword all, or it can be a numeric value, representing one of these protocols or a different one.
	// A protocol name from /etc/protocols is also allowed.
	// A ! argument before the protocol inverts the test.
	// The number zero is equivalent to all.
	// all will match with all protocols and is taken as default when this option is omitted.
	Protocol string `json:"protocol"`

	// TCP flags specification.
	// tcp_flags expects a dict with the two keys flags and flags_set.
	TCPFlags tcpFlags `json:"tcp_flags"`

	// ctstate is a list of the connection states to match in the conntrack module.
	// Possible states are INVALID, NEW, ESTABLISHED, RELATED, UNTRACKED, SNAT, DNAT
	Ctstate []string `json:"ctstate"`

	// Specifies a match to use, that is, an extension module that tests for a specific property.
	// The set of matches make up the condition under which a target is invoked.
	// Matches are evaluated first to last if specified as an array and work in short-circuit fashion,
	// i.e. if one extension yields false, evaluation will stop.
	Match []string `json:"match"`
}

type tcpFlags struct {
	// List of flags you want to examine.
	Flags []string `json:"flags"`
	// Flags to be set.
	FlagsSet []string `json:"flags_set"`
}
