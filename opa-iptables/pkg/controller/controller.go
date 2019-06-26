package controller

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/open-policy-agent/contrib/opa-iptables/pkg/logging"
	"github.com/open-policy-agent/contrib/opa-iptables/pkg/opa"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	ListenAddr  string
	opaEndpoint string
	logger      *logrus.Logger
	opaClient   opa.Client
}

func New(opaEndpoint, hostAddr, hostPort string) *Controller {
	return &Controller{
		opaEndpoint: opaEndpoint,
		logger:      logging.GetLogger(),
		ListenAddr:  hostAddr + ":" + hostPort,
		opaClient:   opa.New(opaEndpoint, ""),
	}
}

func (c *Controller) Run() {
	c.logger.Infof("Controller is running on %s", c.ListenAddr)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	mux := http.NewServeMux()
	mux.HandleFunc("/v0/webhook", c.webhookHandler())

	server := http.Server{
		Addr:         c.ListenAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}

	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			c.logger.Fatal(err)
		}
	}()

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