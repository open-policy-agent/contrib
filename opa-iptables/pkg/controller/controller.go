package controller

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/open-policy-agent/contrib/opa-iptables/pkg/logging"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	ListenAddr  string
	opaEndpoint string
	logger      *logrus.Logger
}

func New(opaEndpoint string) *Controller {
	logger := logging.GetLogger()
	host := "127.0.0.1"
	port := "8080"
	return &Controller{opaEndpoint: opaEndpoint, logger: logger, ListenAddr: host + ":" + port}
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

func (c *Controller) webhookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.logger.Infof("msg=\"Received Request\" req_method=%v req_path=%v\n", r.Method, r.URL)
		if r.Method == http.MethodPost {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				c.logger.Errorf("Error while reading body: %v", err)
				return
			}
			c.logger.Debug(string(body))
		} else {
			res := "This endpoint don't support provided method"
			resStatusCode := http.StatusMethodNotAllowed
			w.Header().Add("Allow", http.MethodPost)
			w.WriteHeader(resStatusCode)
			w.Write([]byte(res))
			c.logger.Errorf("msg=\"Sent Response\" req_method=%v req_path=%v res_bytes=%v res_status=%v\n", r.Method, r.URL, len(res), resStatusCode)
		}
	}
}
