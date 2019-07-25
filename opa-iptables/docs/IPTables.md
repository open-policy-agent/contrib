**Iptables** is used to set up, maintain, and inspect tables of IP packet filter rules in linux kernel. Several different tables may be defined. Each table contains a number of built-in chains and may also contain user-defined chains (Though user-defined chains are *not supported* right now)

Each chain is a list of rules which can match a set of packets. Each rule specifies what to do with a packet that matches. This is called a `target`. Different tables has number of different targets.

To learn more about IPtables you can check out manual page of linux iptables:  https://linux.die.net/man/8/iptables **(Not Newbie Friendly!!!)**

If you are a newbie to IPTable and looking for a tutorial about iptables check out this blog post:
https://www.linode.com/docs/security/firewalls/control-network-traffic-with-iptables/

## How to write rules?

In order to write iptables rules, you have to define rules in JSON format.

Following are the list of parameters for describing the rules:
- [action](#action)
- [chain](#chain)
- [comment](#comment)
- [ctstate](#ctstate)
- [destination](#destination)
- [destination_port](#destination_port)
- [dst_range](#dst_range)
- [in_interface](#in_interface)
- [jump](#jump)
- [match](#match)
- [out_interface](#out_interface)
- [protocol](#protocol)
- [rule_num](#rule_num)
- [source](#source)
- [source_port](#source_port)
- [src_range](#src_range)
- [table](#table)
- [tcp_flags](#tcp_flags)
- [to_destination](#to_destination)
- [to_ports](#to_ports)
- [to_source](#to_source)

---

## action

Whether the rule should be appended at the bottom or inserted at the top.If the rule already exists the chain will not be modified.

Values:
- insert
- append

Default:
- append

Type: `stirng`

## chain

Specify the iptables chain to modify.
The value could be a one of the standard iptables chains, like INPUT, FORWARD, OUTPUT, PREROUTING, POSTROUTING, SECMARK or CONNSECMARK.

Values:
- INPUT
- FORWARD
- OUTPUT
- PREROUTING
- POSTROUTING
- SECMARK
- CONNSECMARK

Default:
- INPUT

Type: `string`

## comment

This specifies a comment that will be added to the rule.

Type: `string`
## ctstate

ctstate is a list of the connection states to match in the conntrack module.
Possible states are INVALID, NEW, ESTABLISHED, RELATED, UNTRACKED, SNAT, DNAT

Values:
- INVALID
- NEW
- ESTABLISHED
- RELATED
- UNTRACKED
- SNAT
- DNAT

Type: `[]string`

## destination

Destination specification.
Address can be either a network IP address (with /mask), or a plain IP address.
The mask can be either a network mask or a plain number, specifying the number of 1's at the left side of the network mask. Thus, a mask of 24 is equivalent to 255.255.255.0. A `!` argument before the address specification inverts the sense of the address.

Type: `string`

## destination_port

Destination port or port range specification. This can either be a service name or a port number. An inclusive range can also be specified, using the format first:last. If the first port is omitted, '0' is assumed; if the last is omitted, '65535' is assumed. If the first port is greater than the second one they will be swapped. This is only valid if the rule also specifies one of the following protocols: tcp, udp, dccp or sctp.

Type: `string`

## dst_range

Specifies the destination IP range to match in the `iprange` module.

Type: `string`

## in_interface

Name of an interface via which a packet was received (only for packets entering the INPUT, FORWARD and PREROUTING chains).
When the ! argument is used before the interface name, the sense is inverted.
If the interface name ends in a +, then any interface which begins with this name will match.
If this option is omitted, any interface name will match.

Type: `string`

## jump

This specifies the target of the rule; i.e., what to do if the packet matches it.
The target can be a user-defined chain (other than the one this rule is in), one of
the special builtin targets which decide the fate of the packet immediately, or an extension (see EXTENSIONS for more at http://ipset.netfilter.org/iptables-extensions.man.html).
If this option is omitted in a rule (and -g is not used), then matching the rule will have no effect
on the packet's fate, but the counters on the rule will be incremented.

Values:

- ACCEPT
- DROP
- REDIRECT
- QUEUE
- RETURN

Type: `string`

## match

Specifies a match to use, that is, an extension module that tests for a specific property.
The set of matches make up the condition under which a target is invoked.
Matches are evaluated first to last if specified as an array and work in short-circuit fashion,
i.e. if one extension yields false, evaluation will stop.

Type: `[]string`

## out_interface

Name of an interface via which a packet is going to be sent (for packets entering the FORWARD, OUTPUT and POSTROUTING chains).
When the ! argument is used before the interface name, the sense is inverted.
If the interface name ends in a +, then any interface which begins with this name will match.
If this option is omitted, any interface name will match.

Type: `string`

## protocol

The protocol of the rule or of the packet to check.
The specified protocol can be one of tcp, udp, udplite, icmp, esp, ah, sctp or the special keyword all, or it can be a numeric value, representing one of these protocols or a different one.
A protocol name from /etc/protocols is also allowed.
A ! argument before the protocol inverts the test.
The number zero is equivalent to all.
all will match with all protocols and is taken as default when this option is omitted.

Values:

- tcp
- udp
- icmp

Type: `string`

## rule_num

Insert the rule as the given rule number.
This works only with action=insert.

Type: `string`

## source

Source Address specification.
Address can be either a network name, a hostname, a network IP address (with /mask), or a plain IP address.
Hostnames will be resolved once only, before the rule is submitted to the kernel.
Please note that specifying any name to be resolved with a remote query such as DNS is a really bad idea.
The mask can be either a network mask or a plain number, specifying the number of 1's at the left side of the network mask.
Thus, a mask of 24 is equivalent to 255.255.255.0. A ! argument before the address specification inverts the sense of the address

Type: `string`

## source_port

Source port or port range specification.
This can either be a service name or a port number.
An inclusive range can also be specified, using the format first:last.
If the first port is omitted, 0 is assumed; if the last is omitted, 65535 is assumed.
If the first port is greater than the second one they will be swapped.

Type: `string`

## src_range

Specifies the source IP range to match in the iprange module.

Type: `string`

## table

This option specifies the packet matching table which the command should operate on.
If the kernel is configured with automatic module loading, an attempt will be made to
load the appropriate module for that table if it is not already there.

Values: 
- filter
- nat
- mangle
- raw
- security

Default: 
- filter

Type: `string`

## tcp_flags

TCP flags specification.
tcp_flags expects a struct with the two keys flags and flags_set.

Type: `object`

Example:

```
"tcp_flags" : {
    "flags": ["ALL"],
    "flags_set": ["ACK","RST","SYN","FIN"]
}
```

## to_destination

This specifies a destination address to use with DNAT.
Without this, the destination address is never altered.

Type: `string`

## to_ports

This specifies a destination port or range of ports to use, without this, the destination port is never altered.
This is only valid if the rule also specifies one of the protocol tcp, udp, dccp or sctp.

Type: `string`

## to_source

This specifies a source address to use with SNAT.
Without this, the source address is never altered.

Type: `string`


