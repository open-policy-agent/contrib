package controller

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/open-policy-agent/contrib/opa-iptables/pkg/logging"
	"github.com/open-policy-agent/contrib/opa-iptables/pkg/opa"
)

func New(config Config) *Controller {
	return &Controller{
		logger:     logging.GetLogger(),
		listenAddr: config.ControllerAddr + ":" + config.ControllerPort,
		opaClient:  opa.New(config.OpaEndpoint, ""),
		w: &watcher{
			watcherInterval: config.WatcherInterval,
			watcherState:    make(map[string]*state),
			watcherDoneCh:   make(chan struct{}, 1),
			logger:          logging.GetLogger(),
		},
		watcherWorkerCount: config.WorkerCount,
		experimental:       config.Experimental,
	}
}

func (c *Controller) Run() {
	c.logger.Infof("Controller is running on %s", c.listenAddr)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	r := mux.NewRouter()
	r.HandleFunc("/v1/iptables/insert", c.insertRuleHandler()).Methods("POST").Queries("q", "")
	r.HandleFunc("/v1/iptables/delete", c.deleteRuleHandler()).Methods("POST").Queries("q", "")
	r.HandleFunc("/v1/iptables/json", c.jsonRuleHandler()).Methods("POST")
	r.HandleFunc("/v1/iptables/list/{table}/{chain}", c.listRulesHandler()).Methods("GET")
	r.HandleFunc("/v1/iptables/list/all", c.listAllRulesHandler()).Methods("GET")

	c.server = http.Server{
		Addr:         c.listenAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
	}

	go c.startController()

	if c.experimental {
		go c.startWatcher()
	}

	<-signalCh
	c.logger.Info("Received SIGINT SIGNAL")

	if c.experimental {
		c.shutdownWatcher()
	}

	c.shutdownController()
}

func (c *Controller) startWatcher() {
	c.newWatcher()
}

func (c *Controller) shutdownWatcher() {
	c.logger.Info("shutting down watcher")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := c.stopWatcher(ctx)
	if err != nil {
		c.logger.Error(err)
	} else {
		c.logger.Info("watcher shutdown successfully")
	}
}

func (c *Controller) startController() {
	err := c.server.ListenAndServe()
	if err != http.ErrServerClosed {
		c.logger.Fatal(err)
	}
}

func (c *Controller) shutdownController() {
	c.logger.Info("shutting down controller")

	ctx, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel1()

	err := c.server.Shutdown(ctx)
	if err != nil {
		if err == context.DeadlineExceeded {
			c.logger.Info("shutdown timeout")
		} else {
			c.logger.Infof("Error while shutting down controller: %s", err)
		}
	} else {
		c.logger.Info("controller shutdown successfully")
	}
}
