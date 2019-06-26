package main

import (
	"flag"

	"github.com/open-policy-agent/contrib/opa-iptables/pkg/controller"
	"github.com/open-policy-agent/contrib/opa-iptables/pkg/logging"
	"github.com/sirupsen/logrus"
)

func main() {
	OpaEndpoint := flag.String("opa-endpoint", "http://127.0.0.1:8181", "endpoint of opa in form of ip:port i.e. 192.33.0.1:8181")
	ControllerAddr := flag.String("controller-host", "0.0.0.0", "controller host")
	ControllerPort := flag.String("controller-port", "8080", "controller port on which it listen on")
	LogFormat := flag.String("log-format", "text", "set log format. i.e. text | json | json-pretty")
	LogLevel := flag.String("log-level", "info", "set log level. i.e. info | debug | error")

	flag.Parse()

	logConfig := logging.Config{
		Format: *LogFormat,
		Level:  *LogLevel,
	}
	logging.SetupLogging(logConfig)

	logger := logging.GetLogger()

	logger.WithFields(logrus.Fields{
		"OPA Endpoint": *OpaEndpoint,
		"Log Format":   *LogFormat,
		"Log Level":    *LogLevel,
	}).Info("Started Controller with following configuration:")

	c := controller.New(*OpaEndpoint,*ControllerAddr,*ControllerPort)
	c.Run()
}
