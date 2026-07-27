// Copyright 2026 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package otlpserverlogger

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc/credentials"

	"github.com/open-policy-agent/opa/util"
	// plugins and version use the /v1/ sub-path; both are in the same go module.
	"github.com/open-policy-agent/opa/v1/plugins"
	opaversion "github.com/open-policy-agent/opa/v1/version"
)

const (
	// PluginName is the name used to register and reference this plugin in OPA config.
	PluginName = "otlp_server_logger"

	// instrumentationScope is the OTel instrumentation scope name — the Go module
	// path of this plugin, per the OpenTelemetry specification.
	instrumentationScope = "github.com/open-policy-agent/contrib/otlp_server_logger"

	// default gRPC port defined in https://opentelemetry.io/docs/specs/otlp/#otlpgrpc-default-port
	defaultGRPCAddress = "localhost:4317"
	// default HTTP port defined in https://opentelemetry.io/docs/specs/otlp/#otlphttp-default-port
	defaultHTTPAddress = "localhost:4318"

	defaultServiceName      = "opa"
	defaultEncryptionScheme = "off"
	defaultLevel            = "info"
)

// Config holds the configuration for the OTLP server logger plugin.
type Config struct {
	Service          string            `json:"service,omitempty"`
	Type             string            `json:"type,omitempty"`
	Address          string            `json:"address,omitempty"`
	ServiceName      string            `json:"service_name,omitempty"`
	EncryptionScheme string            `json:"encryption,omitempty"`
	TLSCertFile      string            `json:"tls_cert_file,omitempty"`
	TLSKeyFile       string            `json:"tls_private_key_file,omitempty"`
	TLSCACertFile    string            `json:"tls_ca_cert_file,omitempty"`
	Level            string            `json:"level,omitempty"`
	Headers          map[string]string `json:"headers,omitempty"`
	Compression      string            `json:"compression,omitempty"`
	ExportTimeoutMs  *int              `json:"export_timeout_ms,omitempty"`
}

// distributedTracingDefaults holds the connection-related subset of OPA's
// distributed_tracing config. Trace-specific fields (sample_percentage,
// batch_span_processor_options) are intentionally omitted.
type distributedTracingDefaults struct {
	Type          string `json:"type,omitempty"`
	Address       string `json:"address,omitempty"`
	ServiceName   string `json:"service_name,omitempty"`
	Encryption    string `json:"encryption,omitempty"`
	TLSCertFile   string `json:"tls_cert_file,omitempty"`
	TLSKeyFile    string `json:"tls_private_key_file,omitempty"`
	TLSCACertFile string `json:"tls_ca_cert_file,omitempty"`
}

// applyDistributedTracingDefaults fills empty connection fields from OPA's
// distributed_tracing config. Explicit plugin config always takes precedence.
func (c *Config) applyDistributedTracingDefaults(raw json.RawMessage) {
	var dt distributedTracingDefaults
	if err := json.Unmarshal(raw, &dt); err != nil {
		return
	}
	if c.Type == "" {
		c.Type = dt.Type
	}
	if c.Address == "" {
		c.Address = dt.Address
	}
	if c.ServiceName == "" {
		c.ServiceName = dt.ServiceName
	}
	if c.EncryptionScheme == "" {
		c.EncryptionScheme = dt.Encryption
	}
	if c.TLSCertFile == "" {
		c.TLSCertFile = dt.TLSCertFile
	}
	if c.TLSKeyFile == "" {
		c.TLSKeyFile = dt.TLSKeyFile
	}
	if c.TLSCACertFile == "" {
		c.TLSCACertFile = dt.TLSCACertFile
	}
}

// inferFromServiceURL derives address and encryption from an OPA service URL.
// https implies tls (or mtls if hasMTLS). Transport type (grpc/http) must be
// set explicitly — port numbers do not reliably identify the OTLP protocol.
func inferFromServiceURL(rawURL string, hasMTLS bool) (address, encryption string) {
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		return
	}
	address = u.Host
	if strings.EqualFold(u.Scheme, "https") {
		if hasMTLS {
			encryption = "mtls"
		} else {
			encryption = "tls"
		}
	} else {
		encryption = "off"
	}
	return
}

// applyServiceDefaults fills empty connection fields from a named OPA service.
// It extracts address, encryption, TLS, and auth headers from the service config.
// Dynamic credential sources (OAuth2, token files) are not supported and must be
// configured explicitly under plugins.otlp_server_logger.
func (c *Config) applyServiceDefaults(m *plugins.Manager) {
	cfg := m.Client(c.Service).Config()

	hasMTLS := cfg.Credentials.ClientTLS != nil
	address, encryption := inferFromServiceURL(cfg.URL, hasMTLS)
	if c.Address == "" {
		c.Address = address
	}
	if c.EncryptionScheme == "" {
		c.EncryptionScheme = encryption
	}

	if c.TLSCACertFile == "" && cfg.TLS != nil && cfg.TLS.CACert != "" {
		c.TLSCACertFile = cfg.TLS.CACert
	}

	if hasMTLS {
		if c.TLSCertFile == "" {
			c.TLSCertFile = cfg.Credentials.ClientTLS.Cert
		}
		if c.TLSKeyFile == "" {
			c.TLSKeyFile = cfg.Credentials.ClientTLS.PrivateKey
		}
	}

	if len(cfg.Headers) > 0 && len(c.Headers) == 0 {
		c.Headers = make(map[string]string, len(cfg.Headers))
		maps.Copy(c.Headers, cfg.Headers)
	}

	if cfg.Credentials.Bearer != nil && cfg.Credentials.Bearer.Token != "" {
		if c.Headers == nil {
			c.Headers = make(map[string]string)
		}
		if _, exists := c.Headers["Authorization"]; !exists {
			c.Headers["Authorization"] = "Bearer " + cfg.Credentials.Bearer.Token
		}
	}
}

// installOTelErrorHandler ensures the global OTel error handler is set at most once,
// preventing repeated global mutation across Stop/Start cycles.
var installOTelErrorHandler sync.Once

func (c *Config) validateAndInjectDefaults() error {
	switch strings.ToLower(c.Type) {
	case "grpc", "http": // OK
	case "":
		return fmt.Errorf("otlp_server_logger.type is required: must be \"grpc\" or \"http\"")
	default:
		return fmt.Errorf("unknown otlp_server_logger.type %q, must be \"grpc\" or \"http\"", c.Type)
	}

	c.Type = strings.ToLower(c.Type)

	if c.Address == "" {
		switch c.Type {
		case "grpc":
			c.Address = defaultGRPCAddress
		case "http":
			c.Address = defaultHTTPAddress
		}
	}

	if strings.Contains(c.Address, "://") {
		return fmt.Errorf("otlp_server_logger.address must not include a URL scheme (got %q); use \"host:port\" format", c.Address)
	}

	if strings.Contains(c.Address, "/") {
		return fmt.Errorf("otlp_server_logger.address must not include a path or trailing slash (got %q); use \"host:port\" format", c.Address)
	}

	if c.ServiceName == "" {
		c.ServiceName = defaultServiceName
	}

	if c.EncryptionScheme == "" {
		c.EncryptionScheme = defaultEncryptionScheme
	}

	c.EncryptionScheme = strings.ToLower(c.EncryptionScheme)

	switch c.EncryptionScheme {
	case "off", "tls", "mtls": // OK
	default:
		return fmt.Errorf("unsupported otlp_server_logger.encryption %q", c.EncryptionScheme)
	}

	if c.EncryptionScheme == "mtls" && (c.TLSCertFile == "" || c.TLSKeyFile == "") {
		return fmt.Errorf("otlp_server_logger.encryption \"mtls\" requires both tls_cert_file and tls_private_key_file to be set")
	}

	if c.Level == "" {
		c.Level = defaultLevel
	}

	if _, err := parseLevel(c.Level); err != nil {
		return err
	}

	c.Level = strings.ToLower(c.Level)

	c.Compression = strings.ToLower(c.Compression)

	switch c.Compression {
	case "", "none", "gzip": // OK
	default:
		return fmt.Errorf("unsupported otlp_server_logger.compression %q, must be \"gzip\" or \"none\"", c.Compression)
	}

	if c.ExportTimeoutMs != nil && *c.ExportTimeoutMs <= 0 {
		return fmt.Errorf("otlp_server_logger.export_timeout_ms must be a positive integer, got %d", *c.ExportTimeoutMs)
	}

	return nil
}

// structuralChange reports whether differences between old and newCfg require the
// exporter to be torn down and rebuilt rather than hot-reloaded.
func structuralChange(old, newCfg Config) bool {
	timeoutChanged := (old.ExportTimeoutMs == nil) != (newCfg.ExportTimeoutMs == nil) ||
		(old.ExportTimeoutMs != nil && newCfg.ExportTimeoutMs != nil &&
			*old.ExportTimeoutMs != *newCfg.ExportTimeoutMs)
	return old.Type != newCfg.Type ||
		old.Address != newCfg.Address ||
		old.ServiceName != newCfg.ServiceName ||
		old.EncryptionScheme != newCfg.EncryptionScheme ||
		old.TLSCertFile != newCfg.TLSCertFile ||
		old.TLSKeyFile != newCfg.TLSKeyFile ||
		old.TLSCACertFile != newCfg.TLSCACertFile ||
		old.Compression != newCfg.Compression ||
		// maps.Equal treats nil and empty map as equal (both len==0); nil→{} is not
		// a meaningful change and correctly does not trigger a restart.
		!maps.Equal(old.Headers, newCfg.Headers) ||
		timeoutChanged
}

// leveledHandler gates log records by level; otelslog.NewHandler ignores
// slog.HandlerOptions.Level, so we override Enabled() with a hot-swappable LevelVar.
type leveledHandler struct {
	slog.Handler
	level *slog.LevelVar
}

func (h leveledHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.level.Level()
}

func (h leveledHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return leveledHandler{h.Handler.WithAttrs(attrs), h.level}
}

func (h leveledHandler) WithGroup(name string) slog.Handler {
	return leveledHandler{h.Handler.WithGroup(name), h.level}
}

// multiHandler fans out log records to multiple slog.Handler implementations.
// Enabled reports true if any handler is enabled for the given level.
type multiHandler []slog.Handler

func (h multiHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for _, handler := range h {
		if handler.Enabled(ctx, l) {
			return true
		}
	}
	return false
}

func (h multiHandler) Handle(ctx context.Context, r slog.Record) error {
	var errs []error
	for _, handler := range h {
		if handler.Enabled(ctx, r.Level) {
			if err := handler.Handle(ctx, r.Clone()); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}

func (h multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make(multiHandler, len(h))
	for i, handler := range h {
		newHandlers[i] = handler.WithAttrs(attrs)
	}
	return newHandlers
}

func (h multiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make(multiHandler, len(h))
	for i, handler := range h {
		newHandlers[i] = handler.WithGroup(name)
	}
	return newHandlers
}

// Factory creates OTLP server logger plugin instances.
type Factory struct{}

// Validate parses and validates the raw plugin configuration bytes.
func (Factory) Validate(m *plugins.Manager, config []byte) (any, error) {
	var cfg Config
	if err := util.Unmarshal(config, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse otlp_server_logger config: %w", err)
	}
	if m != nil {
		if cfg.Service != "" {
			found := false
			for _, svc := range m.Services() {
				if svc == cfg.Service {
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("otlp_server_logger.service %q not found in services", cfg.Service)
			}
			cfg.applyServiceDefaults(m)
		} else if opaConfig := m.GetConfig(); opaConfig != nil && len(opaConfig.DistributedTracing) > 0 {
			cfg.applyDistributedTracingDefaults(opaConfig.DistributedTracing)
		}
	}
	if err := cfg.validateAndInjectDefaults(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// New creates a new Plugin instance from the validated configuration.
func (Factory) New(manager *plugins.Manager, config any) plugins.Plugin {
	return &Plugin{
		manager: manager,
		config:  config.(Config),
	}
}

// Plugin implements plugins.LoggerPlugin, forwarding OPA's internal operational logs
// (the m.logger path) to an OpenTelemetry collector. Note: OPA's console logger —
// used by the decision_logs and status plugins for human-readable stdout output — is
// a separate stream and is not captured by this plugin.
type Plugin struct {
	manager  *plugins.Manager
	config   Config
	handler  slog.Handler
	provider *sdklog.LoggerProvider
	levelVar *slog.LevelVar
	mtx      sync.Mutex
}

// Start initializes the OTLP exporter, logger provider, and slog handler.
func (p *Plugin) Start(ctx context.Context) error {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	return p.startLocked(ctx)
}

// startLocked initializes the exporter and provider. Caller must hold p.mtx.
func (p *Plugin) startLocked(ctx context.Context) error {
	if p.handler != nil {
		return errors.New("otlp server logger already started")
	}

	exporter, err := p.newExporter(ctx)
	if err != nil {
		return fmt.Errorf("create OTLP log exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(p.config.ServiceName),
			semconv.ServiceVersionKey.String(opaversion.Version),
		),
		resource.WithHost(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		// resource.New returns a non-nil partial resource alongside ErrPartialResource
		// (e.g. os.Hostname() failure in a restricted container). Use the partial
		// resource — host.name is informational.
		slog.Warn("otlp_server_logger: resource detection incomplete, proceeding with partial resource", "error", err)
	}

	p.provider = sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
	)

	// Surface SDK-internal export errors (e.g. network failures) via slog.
	// Use sync.Once to avoid repeated global mutation across Stop/Start cycles.
	installOTelErrorHandler.Do(func() {
		otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
			slog.Error("otlp_server_logger: export error", "error", err)
		}))
	})

	level, _ := parseLevel(p.config.Level) // already validated
	p.levelVar = new(slog.LevelVar)
	p.levelVar.Set(level)

	otlpHandler := otelslog.NewHandler(instrumentationScope, otelslog.WithLoggerProvider(p.provider))
	stdoutHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: p.levelVar})

	p.handler = leveledHandler{
		Handler: multiHandler{stdoutHandler, otlpHandler},
		level:   p.levelVar,
	}

	if p.manager != nil {
		p.manager.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateOK})
	}
	return nil
}

// Stop shuts down the logger provider, flushing any pending log records.
func (p *Plugin) Stop(ctx context.Context) {
	// Ensure shutdown is bounded even when the caller provides context.Background().
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.stopLocked(ctx)
}

// stopLocked shuts down the provider and clears plugin state. Caller must hold p.mtx.
func (p *Plugin) stopLocked(ctx context.Context) {
	if p.provider == nil {
		return
	}

	if err := p.provider.Shutdown(ctx); err != nil {
		slog.Error("otlp_server_logger: shutdown failed, pending logs may be lost", "error", err)
	}
	p.provider = nil
	p.handler = nil
	p.levelVar = nil

	if p.manager != nil {
		p.manager.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateNotReady})
	}
}

// Reconfigure updates the plugin configuration.
func (p *Plugin) Reconfigure(ctx context.Context, config any) {
	newConfig := config.(Config)

	p.mtx.Lock()
	defer p.mtx.Unlock()

	oldConfig := p.config
	levelVar := p.levelVar
	p.config = newConfig

	if structuralChange(oldConfig, newConfig) {
		// Bound the shutdown to avoid blocking indefinitely when the exporter
		// is hung (OPA's discovery path passes a no-deadline context).
		stopCtx := ctx
		if _, ok := ctx.Deadline(); !ok {
			var cancel context.CancelFunc
			stopCtx, cancel = context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
		}
		p.stopLocked(stopCtx)
		if err := p.startLocked(ctx); err != nil {
			p.config = oldConfig
			if recoverErr := p.startLocked(ctx); recoverErr != nil { // attempt recovery with previous config
				slog.Error("otlp_server_logger: recovery start failed after config update", "error", recoverErr)
			}
			// Only report StateErr if recovery also failed (plugin is truly stopped).
			// If recovery succeeded, startLocked already set StateOK; don't override it.
			if p.manager != nil && p.handler == nil {
				p.manager.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateErr, Message: err.Error()})
			}
		}
		return
	}

	if oldConfig.Level != newConfig.Level && levelVar != nil {
		// levelVar is nil when the plugin is stopped; skip the update in that case.
		newLevel, err := parseLevel(newConfig.Level)
		if err == nil {
			levelVar.Set(newLevel)
		}
	}
}

func (p *Plugin) Logger() slog.Handler {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	return p.handler
}

func grpcTLSOption(encryptionScheme string, tlsConfig *tls.Config) otlploggrpc.Option {
	if encryptionScheme == "off" {
		return otlploggrpc.WithInsecure()
	}
	return otlploggrpc.WithTLSCredentials(credentials.NewTLS(tlsConfig))
}

func httpTLSOption(encryptionScheme string, tlsConfig *tls.Config) otlploghttp.Option {
	if encryptionScheme == "off" {
		return otlploghttp.WithInsecure()
	}
	return otlploghttp.WithTLSClientConfig(tlsConfig)
}

func (p *Plugin) newExporter(ctx context.Context) (sdklog.Exporter, error) {
	tlsCfg, err := p.buildTLSConfig()
	if err != nil {
		return nil, err
	}

	if p.config.Type == "grpc" {
		opts := []otlploggrpc.Option{
			otlploggrpc.WithEndpoint(p.config.Address),
			grpcTLSOption(p.config.EncryptionScheme, tlsCfg),
		}
		if len(p.config.Headers) > 0 {
			opts = append(opts, otlploggrpc.WithHeaders(p.config.Headers))
		}
		if p.config.Compression == "gzip" {
			opts = append(opts, otlploggrpc.WithCompressor("gzip"))
		}
		if p.config.ExportTimeoutMs != nil {
			opts = append(opts, otlploggrpc.WithTimeout(time.Duration(*p.config.ExportTimeoutMs)*time.Millisecond))
		}
		return otlploggrpc.New(ctx, opts...)
	}
	opts := []otlploghttp.Option{
		otlploghttp.WithEndpoint(p.config.Address),
		httpTLSOption(p.config.EncryptionScheme, tlsCfg),
	}
	if len(p.config.Headers) > 0 {
		opts = append(opts, otlploghttp.WithHeaders(p.config.Headers))
	}
	if p.config.Compression == "gzip" {
		opts = append(opts, otlploghttp.WithCompression(otlploghttp.GzipCompression))
	}
	if p.config.ExportTimeoutMs != nil {
		opts = append(opts, otlploghttp.WithTimeout(time.Duration(*p.config.ExportTimeoutMs)*time.Millisecond))
	}
	return otlploghttp.New(ctx, opts...)
}

func (p *Plugin) buildTLSConfig() (*tls.Config, error) {
	if p.config.EncryptionScheme == "off" {
		return nil, nil
	}

	tlsCfg := &tls.Config{MinVersion: tls.VersionTLS12}

	if p.config.TLSCertFile != "" && p.config.TLSKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(p.config.TLSCertFile, p.config.TLSKeyFile)
		if err != nil {
			return nil, fmt.Errorf("load TLS key pair: %w", err)
		}
		tlsCfg.Certificates = []tls.Certificate{cert}
	}

	if p.config.TLSCACertFile != "" {
		pem, err := os.ReadFile(p.config.TLSCACertFile)
		if err != nil {
			return nil, fmt.Errorf("read CA cert file: %w", err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(pem) {
			return nil, errors.New("no valid certificates found in CA cert file")
		}
		tlsCfg.RootCAs = pool
	}

	return tlsCfg, nil
}

func parseLevel(s string) (slog.Level, error) {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("unknown level %q, must be one of: debug, info, warn, error", s)
	}
}
