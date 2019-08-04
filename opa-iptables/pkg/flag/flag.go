package flag

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/open-policy-agent/contrib/opa-iptables/pkg/iptables"
)

// A Flag represents the state of a flag.
type Flag struct {
	name    string // name as it appear in commnadline
	value   value  // value as set. Each flag must need to statisfy "value" interface in order to set flag value.
	numArgs int    // number of arguments a flag requires
}

// value is the interface to the dynamic value stored in a flag. (The default value is represented as a string.)
// In order to set value of flag, a flag must statisfied value interface
// Set is called once(during parsing), in command line order, for each flag present in FlagSet.
// For instance, the caller could create a flag that turns a comma-separated string into a slice of strings by giving the slice the methods of value; in particular, Set would decompose the comma-separated string into the slice.
type value interface {
	String() string
	Set(string) error
}

// stringValue represents a flag which have a type "string" as a value.
type stringValue string

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) String() string {
	return string(*s)
}

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

// AddStringFlag adds a flag which have a type "string" as a value to a FlagSet.
func (fs *FlagSet) AddStringFlag(p *string, name string, value string, numArgs int) {
	fs.AddFlag(newStringValue(value, p), name, numArgs)
}

// TCPFlags is a struct describes --tcp-flags iptable commandline flag.
type TCPFlags iptables.TcpFlags

// Set value of TCPFlags struct
// provided string have a delimeter '#' which is added while parsing
// actual commandline : --tcp-flags ALL ACK,FIN
// parsed commandline : --tcp-flags="ALL#ACK,FIN"
// TCPFlags.Set("ALL#ACK,FIN")
// TCPFlags.Flags = []string{"ALL"}
// TCPFlags.FlagSet = []string{"ACK","FIN"}
func (t *TCPFlags) Set(val string) error {
	tcp := strings.Split(val, "#")
	if len(tcp) == 2 {
		t.Flags = strings.Split(tcp[0], ",")
		t.FlagsSet = strings.Split(tcp[1], ",")
	} else {
		return fmt.Errorf("invalid value of --tcp-flags")
	}
	return nil
}

func (t *TCPFlags) String() string {
	return fmt.Sprintf("--tcp-flags %v %v", strings.Join(t.Flags, ","), strings.Join(t.FlagsSet, ","))
}

// ErrorHandling defines how FlagSet.Parse behaves if the parse fails.
type errorHandling int

// These constants cause FlagSet.Parse to behave as described if the parse fails.
const (
	ContinueOnError errorHandling = iota // Return a descriptive error.
	ExitOnError                          // Call os.Exit(2).
	PanicOnError                         // Call panic with a descriptive error.
)

var errParse = errors.New("parse error")

type flagError struct {
	name string
}

func (f flagError) Error() string {
	return fmt.Sprintf("flag provided but not defined in FlagSet: -%s", f.name)
}

// descibes error got during parsing flag's arguments
type argumentError struct {
	flagName string
	numArg   int
}

func (a argumentError) Error() string {
	return fmt.Sprintf("flag %s requires (%v) argument", a.flagName, a.numArg)
}

// descibes error got during setting flag's value
type valueError struct {
	flagName  string
	flagValue string
	err       error
}

func (v valueError) Error() string {
	return fmt.Sprintf("invalid value %q for flag -%s: %v", v.flagValue, v.flagName, v.err)
}

// A FlagSet represents a set of defined flags.
// FlagSet uses map for storing a flag. It provides faster lookup while parsing.
type FlagSet struct {
	name          string
	args          []string
	actual        map[string]*Flag
	parsed        bool
	errorHandling errorHandling
}

// NewFlagSet returns a new, empty flag set with the specified name and error handling property.
func NewFlagSet(name string, errorHandling errorHandling) *FlagSet {
	return &FlagSet{name: name, errorHandling: errorHandling}
}

// AddFlag adds a flag with the specified name and number of arguments into FlagSet.
// The type and value of the flag are represented by the first argument, of type "value", which typically holds a user-defined implementation of value interface.
func (fs *FlagSet) AddFlag(value value, name string, numArgs int) {
	flag := &Flag{name, value, numArgs}
	_, alreadythere := fs.actual[name]
	if alreadythere {
		var msg string
		if fs.name == "" {
			msg = fmt.Sprintf("flag redefined: %s", name)
		} else {
			msg = fmt.Sprintf("%s flag redefined: %s", fs.name, name)
		}
		fmt.Println(msg)
		panic(msg) // Happens only if flags are declared with identical names
	}
	if fs.actual == nil {
		fs.actual = make(map[string]*Flag)
	}
	fs.actual[name] = flag
}

func (fs *FlagSet) parse(arguments []string) error {
	fs.parsed = true
	var actualArg []string

	//removes all empty string from arguments
	for _, arg := range arguments {
		if arg != "" {
			actualArg = append(actualArg, arg)
		}
	}

	fs.args = actualArg
	for {
		seen, err := fs.parseOne()
		if seen {
			continue
		}
		if err == nil {
			break
		}
		switch fs.errorHandling {
		case ContinueOnError:
			return err
		case ExitOnError:
			os.Exit(2)
		case PanicOnError:
			panic(err)
		}
	}
	return nil
}

// parseOne parses one commandline arguments at a time.
func (fs *FlagSet) parseOne() (bool, error) {
	if len(fs.args) == 0 {
		return false, nil
	}
	s := fs.args[0]

	if len(s) < 2 || s[0] != '-' {
		return false, fmt.Errorf("%v is not a flag, flag must starts with '-' or '--' and length must be greater than one",s)
	}

	numMinus := 1
	if s[1] == '-' {
		numMinus++
		if len(s) == 2 {
			fs.args = fs.args[1:]
			return false, nil
		}
	}
	name := s[numMinus:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		return false, errParse
	}

	//it's a flag, check if it has any arguments
	fs.args = fs.args[1:]
	var value []string
	flag, isthere := fs.actual[name]
	if !isthere {
		return false, flagError{name}
	}

	// flag must have value, which might be the next arg
	hasArgs := false
	actualValue := ""
	numArg := flag.numArgs
	if len(fs.args)-flag.numArgs >= 0 {
		value, fs.args = fs.args[:numArg], fs.args[numArg:]
		for _, v := range value {
			if len(v) > 0 && v[0] == '-' {
				return false, argumentError{name, numArg}
			}
		}
		hasArgs = true
		actualValue = strings.Join(value, "#")

	}
	if !hasArgs {
		return false, argumentError{name, numArg}
	}
	if err := flag.value.Set(actualValue); err != nil {
		return false, valueError{name, actualValue, err}
	}
	return true, nil
}

// Parse parses arguments list
func (fs *FlagSet) Parse(arguments []string) error {
	if len(arguments) > 1 {
		return fs.parse(arguments[1:])
	}
	return nil
}

// IPTableflagSet represents all possible iptables flag that we support currently.
type IPTableflagSet struct {
	TableFlag        string
	ChainFlag        string
	ProtocolFlag     string
	SourceFlag       string
	DestinationFlag  string
	DportFlag        string
	SportFlag        string
	InInterfaceFlag  string
	OutInterfaceFlag string
	DesRangeFlag     string
	SrcRangeFlag     string
	JumpFlag         string
	MatchFlag        string
	ToPortFlag       string
	CTStateFlag      string
	Comment          string
	LogPrefixFlag    string

	TCPFlag TCPFlags
}

// InitFlagSet Adds user defined Flag into FlagSet.
func (fs *FlagSet) InitFlagSet(tf *IPTableflagSet) {
	fs.AddStringFlag(&tf.TableFlag, "t", "", 1)
	fs.AddStringFlag(&tf.ChainFlag, "A", "", 1)
	fs.AddStringFlag(&tf.ChainFlag, "I", "", 1)
	fs.AddStringFlag(&tf.ProtocolFlag, "p", "", 1)
	fs.AddStringFlag(&tf.SourceFlag, "s", "", 1)
	fs.AddStringFlag(&tf.SourceFlag, "source", "", 1)
	fs.AddStringFlag(&tf.DestinationFlag, "d", "", 1)
	fs.AddStringFlag(&tf.DestinationFlag, "destination", "", 1)
	fs.AddStringFlag(&tf.DportFlag, "dport", "", 1)
	fs.AddStringFlag(&tf.DportFlag, "destination-port", "", 1)
	fs.AddStringFlag(&tf.SportFlag, "sport", "", 1)
	fs.AddStringFlag(&tf.SportFlag, "source-port", "", 1)
	fs.AddStringFlag(&tf.InInterfaceFlag, "i", "", 1)
	fs.AddStringFlag(&tf.InInterfaceFlag, "in-interface", "", 1)
	fs.AddStringFlag(&tf.OutInterfaceFlag, "o", "", 1)
	fs.AddStringFlag(&tf.OutInterfaceFlag, "out-interface", "", 1)
	fs.AddStringFlag(&tf.DesRangeFlag, "dst-range", "", 1)
	fs.AddStringFlag(&tf.SrcRangeFlag, "src-range", "", 1)
	fs.AddStringFlag(&tf.JumpFlag, "j", "", 1)
	fs.AddStringFlag(&tf.MatchFlag, "m", "", 1)
	fs.AddStringFlag(&tf.ToPortFlag, "to-ports", "", 1)
	fs.AddStringFlag(&tf.CTStateFlag, "ctstate", "", 1)
	fs.AddStringFlag(&tf.Comment, "comment", "", 1)
	fs.AddStringFlag(&tf.LogPrefixFlag, "log-prefix", "", 1)
	fs.AddFlag(&tf.TCPFlag, "tcp-flags", 2)
}
