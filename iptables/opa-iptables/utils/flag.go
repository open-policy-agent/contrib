package utils

import (
	"fmt"
	"os"
	"errors"
	"net"
)

// ValidateEndpointFlag validate the value of flag provided through --opa-endpoint "host:port"
// i.e. check for valid "host:port" value
func ValidateEndpointFlag(endpoint string) error {
	host,port,err := net.SplitHostPort(endpoint)
	if err != nil {
		return err
	}

	// net.ParseIp returns nil if host is not a valid textual represention of valid IP Address i.e. invalid IP Address
	ip := net.ParseIP(host)
	if ip == nil {
		return errors.New("Invalid IP Address")
	}

	// opa server listen on port 8181
	if port != "8181" {
		return errors.New("Invalid Port")
	}
	return nil
}

// ValidateDataDirFlag validate the value of flag provided through --watch-data-dir "path to data directory"
// i.e. check for valid directory path
func ValidateDataDirFlag(datadir string) error {
	info,err := os.Stat(datadir)
	if _,ok := err.(*os.PathError) ; ok {
		return fmt.Errorf("%s: no such directory found, provided through --watch-data-dir flag",datadir)
	}

	if !info.IsDir() {
		return errors.New("Invalid Path of Data Directory")
	}

	return nil
}

// ValidatePolicyDirFlag validate the value of flag provided through --watch-policy-dir "path to policy directory"
// i.e. check for valid directory path
func ValidatePolicyDirFlag(policydir string) error {
	info,err := os.Stat(policydir)
	if _,ok := err.(*os.PathError) ; ok {
		return fmt.Errorf("%s: no such directory found, provided through --watch-policy-dir flag",policydir)
	}

	if !info.IsDir() {
		return errors.New("Invalid Path of Policy Directory")
	}

	return nil
}