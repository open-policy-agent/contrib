package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/open-policy-agent/opa/plugins"
	"github.com/open-policy-agent/opa/plugins/logs"
	"github.com/open-policy-agent/opa/runtime"
	"github.com/open-policy-agent/opa/util"
)

type Config struct {
	Stderr bool `json:"stderr"`
}

type Factory struct{}

func (Factory) New(_ *plugins.Manager, config interface{}) plugins.Plugin {
	return &PrintlnLogger{
		config: config.(Config),
	}
}

func (Factory) Validate(_ *plugins.Manager, config []byte) (interface{}, error) {
	parsedConfig := Config{}
	return parsedConfig, util.Unmarshal(config, &parsedConfig)
}

type PrintlnLogger struct {
	config Config
	mtx    sync.Mutex
}

func (p *PrintlnLogger) Start(ctx context.Context) error {
	return nil
}

func (p *PrintlnLogger) Stop(ctx context.Context) {
}

func (p *PrintlnLogger) Reconfigure(ctx context.Context, config interface{}) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.config = config.(Config)
}

func (p *PrintlnLogger) Log(ctx context.Context, event logs.EventV1) error {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	w := os.Stdout
	if p.config.Stderr {
		w = os.Stderr
	}
	fmt.Fprintln(w, event) // ignoring errors!
	return nil
}

func Init() error {
	runtime.RegisterPlugin("println_decision_logger", Factory{})
	return nil
}
