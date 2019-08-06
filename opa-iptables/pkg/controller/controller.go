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

func NewController(opaEndpoint, hostAddr, hostPort string, watcherInterval time.Duration) *Controller {
	return &Controller{
		logger:          logging.GetLogger(),
		listenAddr:      hostAddr + ":" + hostPort,
		opaClient:       opa.New(opaEndpoint, ""),
		w: &watcher{
			watcherInterval: watcherInterval,
			watcherState: make(map[string]*state),
		},
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

	server := http.Server{
		Addr:         c.listenAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
	}

	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			c.logger.Fatal(err)
		}
	}()

	go c.watch()

	<-signalCh
	c.logger.Info("Received SIGINT SIGNAL")
	c.logger.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		if err == context.DeadlineExceeded {
			c.logger.Info("Shutdown Timeout")
		} else {
			c.logger.Infof("Error while shutting down server: %s", err)
		}
	} else {
		c.logger.Info("Server Shutdown Successfully")
	}
}
