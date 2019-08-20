package main

import (
	"flag"
	"fmt"
	"os"

	"runtime"
	"time"

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
	watcherInterval := flag.Duration("watch-interval", 1*time.Minute, "")
	v := flag.Bool("v", false, "show version information")
	workerCount := flag.Int("worker", 3, "number of workers needed for watcher")
	experimental := flag.Bool("experimental", false, "use experimental features")

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

	if *workerCount < 1 || *workerCount > 10 {
		logger.Fatalf(`Provided worker count "%v" is not valid. It must be between 1 and 10.`, *workerCount)
	}

	controllerConfig := controller.Config{
		OpaEndpoint:     *opaEndpoint,
		ControllerAddr:  *controllerAddr,
		ControllerPort:  *controllerPort,
		WatcherInterval: *watcherInterval,
		Experimental:    *experimental,
		WorkerCount:     *workerCount,
	}

	logger.WithFields(logrus.Fields{
		"OPA Endpoint": controllerConfig.OpaEndpoint,
		"Log Format":   logConfig.Format,
		"Log Level":    logConfig.Level,
	}).Info("Started Controller with following configuration:")

	c := controller.New(controllerConfig)
	c.Run()
}

func iptablesExists() bool {
	if _, err := os.Stat("/sbin/iptables"); os.IsNotExist(err) {
		return false
	}
	return true
}
