package main

import (
	"flag"

	"github.com/open-policy-agent/contrib/opa-iptables/pkg/logging"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func main() {
	OpaEndpoint := flag.String("opa-endpoint","http://127.0.0.1:8181","endpoint of opa in form of ip:port i.e. 192.33.0.1:8181")
	LogFormat := flag.String("log-format","json","set log format. i.e. text | json | json-pretty")
	LogLevel := flag.String("log-level","info","set log level. i.e. info | debug | error")

	flag.Parse()

	logConfig := logging.Config{
		Format: *LogFormat,
		Level: *LogLevel,
	}
	logging.SetupLogging(logConfig)

	logger = logging.GetLogger()

	if *OpaEndpoint == "" {
		flag.Usage()
		logger.Fatal("--opa-endpoint is required flags. Please provides values for those flags!")
	}

	logger.WithFields(logrus.Fields{
		"OPA Endpoint":*OpaEndpoint,
		"Log Format":*LogFormat,
		"Log Level": *LogLevel,
	}).Info("Started Controller with following configuration:")
}

func errorExit(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}