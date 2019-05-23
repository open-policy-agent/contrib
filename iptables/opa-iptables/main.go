package main

import (
	"flag"

	"github.com/contrib/iptables/opa-iptables/logging"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func main() {
	OpaEndpoint := flag.String("opa-endpoint","127.0.0.1:8181","endpoint of opa in form of ip:port i.e. 192.33.0.1:8181")
	DataDir := flag.String("watch-data-dir","","Relative or Absolute path to the data directory for watching any changes")
	PolicyDir := flag.String("watch-policy-dir","","Relative or Absolute path to the policy directory for watching any changes")
	LogFormat := flag.String("log-format","json","set log format. i.e. text | json | json-pretty")
	LogLevel := flag.String("log-level","info","set log level. i.e. info | debug | error")

	flag.Parse()

	logConfig := logging.Config{
		Format: *LogFormat,
		Level: *LogLevel,
	}
	logging.SetupLogging(logConfig)

	logger = logging.GetLogger()

	if *OpaEndpoint == "" || *DataDir == "" || *PolicyDir == "" {
		flag.Usage()
		logger.Fatal("--opa-endpoint | --watch-data-dir | --watch-policy-dir are required flags. Please provides values for those flags!")
	}
	
	err := validateEndpointFlag(*OpaEndpoint)
	errorExit(err)

	err = validateDataDirFlag(*DataDir)
	errorExit(err)

	err = validatePolicyDirFlag(*PolicyDir)
	errorExit(err)

	logger.WithFields(logrus.Fields{
		"OPA Endpoint":*OpaEndpoint,
		"Data Directory":*DataDir,
		"Policy Directory":*PolicyDir,
		"Log Format":*LogFormat,
		"Log Level": *LogLevel,
	}).Info("Started Controller with following configuration:")
}

func errorExit(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}