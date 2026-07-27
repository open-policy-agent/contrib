// Copyright 2026 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package otlpserverlogger

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io"
	"log/slog"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"
)

func TestValidateAndInjectDefaults(t *testing.T) {
	intPtr := func(v int) *int { return &v }

	tests := []struct {
		name            string
		input           Config
		wantAddress     string
		wantService     string
		wantEncrypt     string
		wantLevel       string
		wantType        string
		wantCompression string
		wantErr         bool
	}{
		{
			name:        "grpc defaults injected",
			input:       Config{Type: "grpc"},
			wantAddress: defaultGRPCAddress,
			wantService: defaultServiceName,
			wantEncrypt: defaultEncryptionScheme,
			wantLevel:   defaultLevel,
		},
		{
			name:        "http defaults injected",
			input:       Config{Type: "http"},
			wantAddress: defaultHTTPAddress,
			wantService: defaultServiceName,
			wantEncrypt: defaultEncryptionScheme,
			wantLevel:   defaultLevel,
		},
		{
			name: "explicit values preserved",
			input: Config{
				Type:             "grpc",
				Address:          "collector.example.com:4317",
				ServiceName:      "my-opa",
				EncryptionScheme: "tls",
				Level:            "debug",
			},
			wantAddress: "collector.example.com:4317",
			wantService: "my-opa",
			wantEncrypt: "tls",
			wantLevel:   "debug",
		},
		{
			name:    "unknown type rejected",
			input:   Config{Type: "otlp/udp"},
			wantErr: true,
		},
		{
			name:    "empty type rejected",
			input:   Config{},
			wantErr: true,
		},
		{
			name:    "unsupported encryption rejected",
			input:   Config{Type: "grpc", EncryptionScheme: "starttls"},
			wantErr: true,
		},
		{
			name:        "type matching is case-insensitive",
			input:       Config{Type: "GRPC"},
			wantAddress: defaultGRPCAddress,
			wantService: defaultServiceName,
			wantEncrypt: defaultEncryptionScheme,
			wantLevel:   defaultLevel,
		},
		{
			name:    "invalid level rejected",
			input:   Config{Type: "grpc", Level: "verbose"},
			wantErr: true,
		},
		{
			name:    "address with scheme rejected",
			input:   Config{Type: "grpc", Address: "http://localhost:4317"},
			wantErr: true,
		},
		{
			name:    "address with trailing slash rejected",
			input:   Config{Type: "grpc", Address: "localhost:4317/"},
			wantErr: true,
		},
		{
			name:    "address with path rejected",
			input:   Config{Type: "grpc", Address: "localhost:4317/v1/logs"},
			wantErr: true,
		},
		{
			name:    "mtls without cert files rejected",
			input:   Config{Type: "grpc", EncryptionScheme: "mtls"},
			wantErr: true,
		},
		{
			name:        "level is normalized to lowercase",
			input:       Config{Type: "grpc", Level: "INFO"},
			wantAddress: defaultGRPCAddress,
			wantService: defaultServiceName,
			wantEncrypt: defaultEncryptionScheme,
			wantLevel:   "info",
		},
		{
			name:    "invalid compression rejected",
			input:   Config{Type: "grpc", Compression: "brotli"},
			wantErr: true,
		},
		{
			name:            "gzip compression accepted and normalized",
			input:           Config{Type: "grpc", Compression: "gzip"},
			wantAddress:     defaultGRPCAddress,
			wantService:     defaultServiceName,
			wantEncrypt:     defaultEncryptionScheme,
			wantLevel:       defaultLevel,
			wantCompression: "gzip",
		},
		{
			name:            "GZIP compression normalized to lowercase",
			input:           Config{Type: "grpc", Compression: "GZIP"},
			wantAddress:     defaultGRPCAddress,
			wantService:     defaultServiceName,
			wantEncrypt:     defaultEncryptionScheme,
			wantLevel:       defaultLevel,
			wantCompression: "gzip",
		},
		{
			name:            "NONE compression normalized to lowercase",
			input:           Config{Type: "grpc", Compression: "NONE"},
			wantAddress:     defaultGRPCAddress,
			wantService:     defaultServiceName,
			wantEncrypt:     defaultEncryptionScheme,
			wantLevel:       defaultLevel,
			wantCompression: "none",
		},
		{
			name:        "type GRPC normalized to lowercase",
			input:       Config{Type: "GRPC"},
			wantType:    "grpc",
			wantAddress: defaultGRPCAddress,
			wantService: defaultServiceName,
			wantEncrypt: defaultEncryptionScheme,
			wantLevel:   defaultLevel,
		},
		{
			name:        "type HTTP normalized to lowercase",
			input:       Config{Type: "HTTP"},
			wantType:    "http",
			wantAddress: defaultHTTPAddress,
			wantService: defaultServiceName,
			wantEncrypt: defaultEncryptionScheme,
			wantLevel:   defaultLevel,
		},
		{
			name:    "zero export_timeout_ms rejected",
			input:   Config{Type: "grpc", ExportTimeoutMs: intPtr(0)},
			wantErr: true,
		},
		{
			name:    "negative export_timeout_ms rejected",
			input:   Config{Type: "grpc", ExportTimeoutMs: intPtr(-1)},
			wantErr: true,
		},
		{
			name:        "positive export_timeout_ms accepted",
			input:       Config{Type: "grpc", ExportTimeoutMs: intPtr(5000)},
			wantAddress: defaultGRPCAddress,
			wantService: defaultServiceName,
			wantEncrypt: defaultEncryptionScheme,
			wantLevel:   defaultLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.input
			err := cfg.validateAndInjectDefaults()

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if cfg.Address != tt.wantAddress {
				t.Errorf("Address: got %q, want %q", cfg.Address, tt.wantAddress)
			}
			if cfg.ServiceName != tt.wantService {
				t.Errorf("ServiceName: got %q, want %q", cfg.ServiceName, tt.wantService)
			}
			if cfg.EncryptionScheme != tt.wantEncrypt {
				t.Errorf("EncryptionScheme: got %q, want %q", cfg.EncryptionScheme, tt.wantEncrypt)
			}
			if cfg.Level != tt.wantLevel {
				t.Errorf("Level: got %q, want %q", cfg.Level, tt.wantLevel)
			}
			if tt.wantType != "" && cfg.Type != tt.wantType {
				t.Errorf("Type: got %q, want %q", cfg.Type, tt.wantType)
			}
			if cfg.Compression != tt.wantCompression {
				t.Errorf("Compression: got %q, want %q", cfg.Compression, tt.wantCompression)
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input   string
		want    slog.Level
		wantErr bool
	}{
		{input: "debug", want: slog.LevelDebug},
		{input: "info", want: slog.LevelInfo},
		{input: "warn", want: slog.LevelWarn},
		{input: "error", want: slog.LevelError},
		{input: "DEBUG", want: slog.LevelDebug},
		{input: "INFO", want: slog.LevelInfo},
		{input: "WARN", want: slog.LevelWarn},
		{input: "ERROR", want: slog.LevelError},
		{input: "verbose", wantErr: true},
		{input: "", wantErr: true},
		{input: "trace", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseLevel(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseLevel(%q): expected error, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("parseLevel(%q): unexpected error: %v", tt.input, err)
				return
			}
			if got != tt.want {
				t.Errorf("parseLevel(%q): got %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestStopLockedWhenNeverStarted(t *testing.T) {
	// stopLocked on a zero-value Plugin must not panic (provider == nil guard).
	p := &Plugin{}
	p.stopLocked(context.Background()) // must not panic
}

func TestPluginLoggerNilBeforeStart(t *testing.T) {
	p := &Plugin{}
	if p.Logger() != nil {
		t.Error("expected Logger() to return nil before Start()")
	}
}

func TestPluginAlreadyStartedGuard(t *testing.T) {
	p := &Plugin{}
	p.handler = slog.NewTextHandler(io.Discard, nil)

	err := p.Start(context.Background())
	if err == nil {
		t.Fatal("expected error when starting already-started plugin, got nil")
	}
	if !strings.Contains(err.Error(), "already started") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestPluginReconfigureLevelOnly(t *testing.T) {
	lv := new(slog.LevelVar)
	lv.Set(slog.LevelInfo)

	p := &Plugin{
		config:   Config{Type: "grpc", Address: "localhost:4317", ServiceName: "opa", EncryptionScheme: "off", Level: "info"},
		handler:  slog.NewTextHandler(io.Discard, nil),
		levelVar: lv,
	}

	newCfg := p.config
	newCfg.Level = "debug"
	p.Reconfigure(context.Background(), newCfg)

	if lv.Level() != slog.LevelDebug {
		t.Errorf("expected level Debug after reconfigure, got %v", lv.Level())
	}
}

func TestReconfigureStructuralChangeDetection(t *testing.T) {
	base := Config{
		Type:             "grpc",
		Address:          "localhost:4317",
		ServiceName:      "opa",
		EncryptionScheme: "off",
		Level:            "info",
	}

	intPtr := func(v int) *int { return &v }

	cases := []struct {
		name   string
		mutate func(*Config)
		want   bool
	}{
		{"Type change", func(c *Config) { c.Type = "http" }, true},
		{"Address change", func(c *Config) { c.Address = "collector:4317" }, true},
		{"ServiceName change", func(c *Config) { c.ServiceName = "my-opa" }, true},
		{"EncryptionScheme change", func(c *Config) { c.EncryptionScheme = "tls" }, true},
		{"TLSCertFile change", func(c *Config) { c.TLSCertFile = "/a/cert.pem" }, true},
		{"TLSKeyFile change", func(c *Config) { c.TLSKeyFile = "/a/key.pem" }, true},
		{"TLSCACertFile change", func(c *Config) { c.TLSCACertFile = "/a/ca.pem" }, true},
		{"Compression change", func(c *Config) { c.Compression = "gzip" }, true},
		{"ExportTimeoutMs set", func(c *Config) { c.ExportTimeoutMs = intPtr(5000) }, true},
		{"Level change only — not structural", func(c *Config) { c.Level = "debug" }, false},
		{"Headers change — structural (SDK snapshots headers at creation)", func(c *Config) { c.Headers = map[string]string{"k": "v"} }, true},
		{"No change", func(c *Config) {}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			newCfg := base
			tc.mutate(&newCfg)
			// Calls the package-level structuralChange function — not a local copy.
			got := structuralChange(base, newCfg)
			if got != tc.want {
				t.Errorf("structuralChange: got %v, want %v", got, tc.want)
			}
		})
	}

	t.Run("ExportTimeoutMs pointer semantics", func(t *testing.T) {
		type ptrCase struct {
			name string
			old  *int
			new_ *int
			want bool
		}
		ptrCases := []ptrCase{
			{"both nil — not structural", nil, nil, false},
			{"nil to non-nil — structural", nil, intPtr(5000), true},
			{"non-nil to nil — structural", intPtr(5000), nil, true},
			{"same non-nil value — not structural", intPtr(5000), intPtr(5000), false},
			{"different non-nil values — structural", intPtr(5000), intPtr(3000), true},
		}
		for _, pc := range ptrCases {
			t.Run(pc.name, func(t *testing.T) {
				oldCfg := base
				oldCfg.ExportTimeoutMs = pc.old
				newCfg := base
				newCfg.ExportTimeoutMs = pc.new_
				got := structuralChange(oldCfg, newCfg)
				if got != pc.want {
					t.Errorf("structuralChange=%v, want %v", got, pc.want)
				}
			})
		}
	})
}

func TestLeveledHandlerEnabled(t *testing.T) {
	lv := new(slog.LevelVar)
	lv.Set(slog.LevelWarn)
	h := leveledHandler{slog.NewTextHandler(io.Discard, nil), lv}

	if h.Enabled(context.Background(), slog.LevelDebug) {
		t.Error("expected Debug to be disabled at Warn level")
	}
	if h.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("expected Info to be disabled at Warn level")
	}
	if !h.Enabled(context.Background(), slog.LevelWarn) {
		t.Error("expected Warn to be enabled at Warn level")
	}
	if !h.Enabled(context.Background(), slog.LevelError) {
		t.Error("expected Error to be enabled at Warn level")
	}
}

func TestLeveledHandlerWithAttrs(t *testing.T) {
	lv := new(slog.LevelVar)
	lv.Set(slog.LevelWarn)
	inner := slog.NewTextHandler(io.Discard, nil)
	h := leveledHandler{inner, lv}

	h2 := h.WithAttrs([]slog.Attr{slog.String("k", "v")})
	lh, ok := h2.(leveledHandler)
	if !ok {
		t.Fatalf("WithAttrs returned %T, want leveledHandler", h2)
	}
	if !lh.Enabled(context.Background(), slog.LevelWarn) {
		t.Error("WithAttrs result lost level filtering")
	}
	if lh.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("WithAttrs result lost level filtering")
	}
}

func TestLeveledHandlerWithGroup(t *testing.T) {
	lv := new(slog.LevelVar)
	lv.Set(slog.LevelError)
	inner := slog.NewTextHandler(io.Discard, nil)
	h := leveledHandler{inner, lv}

	h2 := h.WithGroup("mygroup")
	lh, ok := h2.(leveledHandler)
	if !ok {
		t.Fatalf("WithGroup returned %T, want leveledHandler", h2)
	}
	if lh.Enabled(context.Background(), slog.LevelWarn) {
		t.Error("WithGroup result lost level filtering")
	}
	if !lh.Enabled(context.Background(), slog.LevelError) {
		t.Error("WithGroup result lost level filtering")
	}
}

func generateTestCertFiles(t *testing.T) (certFile, keyFile, caFile string) {
	t.Helper()

	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate CA key: %v", err)
	}
	caTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "test-ca"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("create CA cert: %v", err)
	}

	caF, err := os.CreateTemp(t.TempDir(), "ca-*.pem")
	if err != nil {
		t.Fatalf("create CA file: %v", err)
	}
	if err := pem.Encode(caF, &pem.Block{Type: "CERTIFICATE", Bytes: caDER}); err != nil {
		t.Fatalf("encode CA cert: %v", err)
	}
	caF.Close()
	caFile = caF.Name()

	certF, err := os.CreateTemp(t.TempDir(), "cert-*.pem")
	if err != nil {
		t.Fatalf("create cert file: %v", err)
	}
	if err := pem.Encode(certF, &pem.Block{Type: "CERTIFICATE", Bytes: caDER}); err != nil {
		t.Fatalf("encode leaf cert: %v", err)
	}
	certF.Close()
	certFile = certF.Name()

	keyDER, err := x509.MarshalECPrivateKey(caKey)
	if err != nil {
		t.Fatalf("marshal key: %v", err)
	}
	keyF, err := os.CreateTemp(t.TempDir(), "key-*.pem")
	if err != nil {
		t.Fatalf("create key file: %v", err)
	}
	if err := pem.Encode(keyF, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}); err != nil {
		t.Fatalf("encode key: %v", err)
	}
	keyF.Close()
	keyFile = keyF.Name()

	return certFile, keyFile, caFile
}

func TestBuildTLSConfig(t *testing.T) {
	certFile, keyFile, caFile := generateTestCertFiles(t)

	garbageF, err := os.CreateTemp(t.TempDir(), "garbage-*.pem")
	if err != nil {
		t.Fatalf("create garbage file: %v", err)
	}
	if _, err := garbageF.WriteString("this is not a PEM block"); err != nil {
		t.Fatalf("write garbage file: %v", err)
	}
	garbageF.Close()

	tests := []struct {
		name        string
		config      Config
		wantNil     bool
		wantErr     bool
		errContains string
		wantCerts   int
		wantRootCAs bool
		wantMinVer  uint16
	}{
		{
			name:    "off returns nil config",
			config:  Config{EncryptionScheme: "off"},
			wantNil: true,
		},
		{
			name:       "tls no files returns bare config",
			config:     Config{EncryptionScheme: "tls"},
			wantMinVer: tls.VersionTLS12,
			wantCerts:  0,
		},
		{
			name: "tls valid cert pair",
			config: Config{
				EncryptionScheme: "tls",
				TLSCertFile:      certFile,
				TLSKeyFile:       keyFile,
			},
			wantMinVer: tls.VersionTLS12,
			wantCerts:  1,
		},
		{
			name: "tls nonexistent cert file errors",
			config: Config{
				EncryptionScheme: "tls",
				TLSCertFile:      "/does/not/exist.crt",
				TLSKeyFile:       "/does/not/exist.key",
			},
			wantErr:     true,
			errContains: "load TLS key pair",
		},
		{
			name: "tls nonexistent CA file errors",
			config: Config{
				EncryptionScheme: "tls",
				TLSCACertFile:    "/does/not/exist-ca.pem",
			},
			wantErr:     true,
			errContains: "read CA cert file",
		},
		{
			name: "tls CA file with no valid PEM certs errors",
			config: Config{
				EncryptionScheme: "tls",
				TLSCACertFile:    garbageF.Name(),
			},
			wantErr:     true,
			errContains: "no valid certificates found",
		},
		{
			name: "tls valid CA cert file sets RootCAs",
			config: Config{
				EncryptionScheme: "tls",
				TLSCACertFile:    caFile,
			},
			wantMinVer:  tls.VersionTLS12,
			wantRootCAs: true,
		},
		{
			name: "mtls cert pair and CA cert",
			config: Config{
				EncryptionScheme: "mtls",
				TLSCertFile:      certFile,
				TLSKeyFile:       keyFile,
				TLSCACertFile:    caFile,
			},
			wantMinVer:  tls.VersionTLS12,
			wantCerts:   1,
			wantRootCAs: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Plugin{config: tt.config}
			got, err := p.buildTLSConfig()

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.wantNil {
				if got != nil {
					t.Errorf("expected nil *tls.Config, got non-nil")
				}
				return
			}
			if got == nil {
				t.Fatal("expected non-nil *tls.Config, got nil")
			}
			if tt.wantMinVer != 0 && got.MinVersion != tt.wantMinVer {
				t.Errorf("MinVersion: got %d, want %d", got.MinVersion, tt.wantMinVer)
			}
			if len(got.Certificates) != tt.wantCerts {
				t.Errorf("len(Certificates): got %d, want %d", len(got.Certificates), tt.wantCerts)
			}
			if tt.wantRootCAs && got.RootCAs == nil {
				t.Error("expected non-nil RootCAs, got nil")
			}
			if !tt.wantRootCAs && got.RootCAs != nil {
				t.Error("expected nil RootCAs, got non-nil")
			}
		})
	}
}

func TestGRPCTLSOption(t *testing.T) {
	tlsCfg := &tls.Config{MinVersion: tls.VersionTLS12}

	insecureOpt := grpcTLSOption("off", nil)
	if insecureOpt == nil {
		t.Error("expected non-nil option for off")
	}

	tlsOpt := grpcTLSOption("tls", tlsCfg)
	if tlsOpt == nil {
		t.Error("expected non-nil option for tls")
	}
}

func TestHTTPTLSOption(t *testing.T) {
	tlsCfg := &tls.Config{MinVersion: tls.VersionTLS12}

	insecureOpt := httpTLSOption("off", nil)
	if insecureOpt == nil {
		t.Error("expected non-nil option for off")
	}

	tlsOpt := httpTLSOption("tls", tlsCfg)
	if tlsOpt == nil {
		t.Error("expected non-nil option for tls")
	}
}

func TestMultiHandlerEnabled(t *testing.T) {
	enabledHandler := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug})
	disabledHandler := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})

	h := multiHandler{disabledHandler, enabledHandler}
	if !h.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("expected Enabled=true when at least one handler is enabled")
	}

	hAllDisabled := multiHandler{disabledHandler, disabledHandler}
	if hAllDisabled.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("expected Enabled=false when all handlers are disabled")
	}
}

func TestMultiHandlerWithAttrs(t *testing.T) {
	var buf1, buf2 strings.Builder
	h1 := slog.NewTextHandler(&buf1, nil)
	h2 := slog.NewTextHandler(&buf2, nil)

	mh := multiHandler{h1, h2}
	mh2 := mh.WithAttrs([]slog.Attr{slog.String("k", "v")})
	if _, ok := mh2.(multiHandler); !ok {
		t.Fatalf("WithAttrs returned %T, want multiHandler", mh2)
	}
}

func TestMultiHandlerWithGroup(t *testing.T) {
	var buf1, buf2 strings.Builder
	h1 := slog.NewTextHandler(&buf1, nil)
	h2 := slog.NewTextHandler(&buf2, nil)

	mh := multiHandler{h1, h2}
	mh2 := mh.WithGroup("g")
	if _, ok := mh2.(multiHandler); !ok {
		t.Fatalf("WithGroup returned %T, want multiHandler", mh2)
	}
}

// erroringHandler is a slog.Handler that always returns an error from Handle.
type erroringHandler struct {
	err        error
	lastRecord slog.Record
}

func (h *erroringHandler) Enabled(_ context.Context, _ slog.Level) bool { return true }
func (h *erroringHandler) Handle(_ context.Context, r slog.Record) error {
	h.lastRecord = r
	return h.err
}
func (h *erroringHandler) WithAttrs(_ []slog.Attr) slog.Handler { return h }
func (h *erroringHandler) WithGroup(_ string) slog.Handler      { return h }

func TestMultiHandlerHandle(t *testing.T) {
	handlerErr := errors.New("handler error")
	eh := &erroringHandler{err: handlerErr}
	var buf strings.Builder
	okHandler := slog.NewTextHandler(&buf, nil)

	mh := multiHandler{eh, okHandler}
	r := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)

	err := mh.Handle(context.Background(), r)

	// Error from the first handler must be returned.
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if !errors.Is(err, handlerErr) {
		t.Errorf("got error %v, want it to wrap %v", err, handlerErr)
	}
	// Second handler must still be called even when first failed.
	if buf.Len() == 0 {
		t.Error("second handler was not called after first returned an error")
	}
	// Each handler receives an independent clone of the record.
	if eh.lastRecord.Message != "test message" {
		t.Errorf("first handler received record with message %q, want %q", eh.lastRecord.Message, "test message")
	}
}

func BenchmarkValidateAndInjectDefaults(b *testing.B) {
	cfg := Config{Type: "grpc", Level: "info", Compression: "gzip"}
	b.ReportAllocs()
	for range b.N {
		c := cfg
		_ = c.validateAndInjectDefaults()
	}
}

func TestInferFromServiceURL(t *testing.T) {
	tests := []struct {
		name           string
		rawURL         string
		hasMTLS        bool
		wantAddress    string
		wantEncryption string
	}{
		{
			name:           "http scheme gives off encryption",
			rawURL:         "http://otel-collector:4317",
			wantAddress:    "otel-collector:4317",
			wantEncryption: "off",
		},
		{
			name:           "https scheme gives tls",
			rawURL:         "https://otel-collector:4317",
			wantAddress:    "otel-collector:4317",
			wantEncryption: "tls",
		},
		{
			name:           "https with client TLS gives mtls",
			rawURL:         "https://otel-collector:4317",
			hasMTLS:        true,
			wantAddress:    "otel-collector:4317",
			wantEncryption: "mtls",
		},
		{
			name:           "non-standard port is preserved as-is",
			rawURL:         "http://otel-collector:9999",
			wantAddress:    "otel-collector:9999",
			wantEncryption: "off",
		},
		{
			name:           "uppercase HTTPS scheme is case-insensitive",
			rawURL:         "HTTPS://otel-collector:4317",
			wantAddress:    "otel-collector:4317",
			wantEncryption: "tls",
		},
		{
			name:           "URL with path strips path from address",
			rawURL:         "https://otel-collector:4317/v1/logs",
			wantAddress:    "otel-collector:4317",
			wantEncryption: "tls",
		},
		{
			name:           "scheme-only URL with no host returns empty",
			rawURL:         "https://",
			wantAddress:    "",
			wantEncryption: "",
		},
		{
			name:           "host without port is preserved as-is",
			rawURL:         "https://otel-collector",
			wantAddress:    "otel-collector",
			wantEncryption: "tls",
		},
		{
			name:           "IPv6 address is preserved",
			rawURL:         "http://[::1]:4317",
			wantAddress:    "[::1]:4317",
			wantEncryption: "off",
		},
		{
			name:           "http scheme ignores hasMTLS",
			rawURL:         "http://otel-collector:4317",
			hasMTLS:        true,
			wantAddress:    "otel-collector:4317",
			wantEncryption: "off",
		},
		{
			name: "empty URL returns empty fields",
		},
		{
			name:   "invalid URL returns empty fields",
			rawURL: "://bad",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddress, gotEncryption := inferFromServiceURL(tt.rawURL, tt.hasMTLS)
			if gotAddress != tt.wantAddress {
				t.Errorf("address: got %q, want %q", gotAddress, tt.wantAddress)
			}
			if gotEncryption != tt.wantEncryption {
				t.Errorf("encryption: got %q, want %q", gotEncryption, tt.wantEncryption)
			}
		})
	}
}

func BenchmarkParseLevel(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_, _ = parseLevel("info")
	}
}
