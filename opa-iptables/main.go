package main

import (
	"time"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/open-policy-agent/contrib/opa-iptables/pkg/controller"
	"github.com/open-policy-agent/contrib/opa-iptables/pkg/logging"
	"github.com/open-policy-agent/contrib/opa-iptables/pkg/version"
	"github.com/sirupsen/logrus"
)

func main() {
	opaEndpoint := flag.String("opa-endpoint", "http://127.0.0.1:8181", "endpoint of opa in form of http://ip:port i.e. http://192.33.0.1:8181")
	controllerAddr := flag.String("controller-host", "0.0.0.0", "controller host")
	// setting default port value to some high port to prevent accidentally block this port in IPTable rules
	controllerPort := flag.String("controller-port", "33455", "controller port on which it listen on")
	logFormat := flag.String("log-format", "text", "set log format. i.e. text | json | json-pretty")
	logLevel := flag.String("log-level", "info", "set log level. i.e. info | debug | error")
	watcherInterval := flag.Duration("watch-interval",1*time.Minute,"")
	v := flag.Bool("v", false, "show version information")

	flag.Parse()

	if *v {
		fmt.Printf("Version= %v\nCommit= %v\n", version.Version, version.Commit)
		os.Exit(0)
	}

	logConfig := logging.Config{
		Format: *logFormat,
		Level:  *logLevel,
	}
	logging.SetupLogging(logConfig)

	logger := logging.GetLogger()

	if runtime.GOOS != "linux" {
		logger.Errorln("\"iptables\" utility is only supported on Linux kernel. It's seems like that you are not running Linux kernel.")
		os.Exit(1)
	}

	if !iptablesExists() {
		logger.Error("command \"iptables\" not found at path \"/sbin/iptables\".")
		fmt.Println(installationHelp)
		os.Exit(1)
	}

	logger.WithFields(logrus.Fields{
		"OPA Endpoint": *opaEndpoint,
		"Log Format":   *logFormat,
		"Log Level":    *logLevel,
	}).Info("Started Controller with following configuration:")

	c := controller.NewController(*opaEndpoint, *controllerAddr, *controllerPort, *watcherInterval)
	c.Run()
}

func iptablesExists() bool {
	if _, err := os.Stat("/sbin/iptables"); os.IsNotExist(err) {
		return false
	}
	return true
}