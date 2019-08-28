package controller

import (
	"net/http"
	"sync"
	"time"

	"github.com/open-policy-agent/contrib/opa-iptables/pkg/opa"
	"github.com/sirupsen/logrus"
)

type Config struct {
	OpaEndpoint     string
	ControllerAddr  string
	ControllerPort  string
	WatcherInterval time.Duration
	WatcherFlag     bool
	WorkerCount     int
}

// Controller is a struct which is used for storing server related data.
// It contains logger for centralize logging, opaClient for accessing OPA REST API, and
// watcher for watching any state changes in ruleset of registred state stored in
// watcherstate map.
type Controller struct {
	listenAddr         string
	server             http.Server
	logger             *logrus.Logger
	opaClient          opa.Client
	w                  *watcher
	watcherWorkerCount int
	watcher            bool
}

// state is used for storing nessecarry information for doing repeated query for checking
// "_id" field in ruleset. This state is stored in a watcherState map using "queryPath" as
// a key and "state" as a value.
type state struct {
	id        string
	payload   payload
	queryPath string
}

type payload struct {
	Input interface{} `json:"input"`
}

// requerst represents opa-iptables insert/delete rules request API.
// i.e. http://localhost:33455/v1/iptables/insert?q=iptables/webserver_rules
// request body:
// 	{
//		"input" : {}
//  }
//  here, "queryPath" will be : iptables/webserver_rules
//        "payload"  will be : request body
type request struct {
	queryPath string
	p         payload
}

// watcher is used for storing state and checking and updating any state changes.
// watcher checks any changes to watcherState at every "watcherInterval" time duration.
type watcher struct {
	watcherInterval time.Duration
	watcherDoneCh   chan struct{}
	logger          *logrus.Logger

	mu           sync.RWMutex // guard the following fields
	watcherState map[string]*state
}
